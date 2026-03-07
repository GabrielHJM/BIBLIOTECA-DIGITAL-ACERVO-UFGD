package main

// @title Biblioteca Digital API
// @version 1.3.0
// @description API para o sistema de Biblioteca Digital.
// @termsOfService http://swagger.io/terms/

// @contact.name Gabriel
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"
	"biblioteca-digital-api/internal/handler/middleware"
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/ai"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"
	"sync"

	_ "biblioteca-digital-api/docs"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

var (
	startTime = time.Now()
	syncMu    sync.Mutex
)

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)

	// Otimização do Pool de Conexões do Banco de Dados
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	defer db.Close()

	// Core Migrations
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS usuarios (
			id SERIAL PRIMARY KEY,
			nome TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			senha TEXT NOT NULL,
			tipo INTEGER DEFAULT 1,
			foto_url TEXT,
			meta_paginas_semana INTEGER DEFAULT 100,
			data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS materiais (
			id SERIAL PRIMARY KEY,
			titulo TEXT NOT NULL,
			autor TEXT NOT NULL,
			isbn TEXT,
			categoria TEXT NOT NULL,
			ano_publicacao INTEGER,
			descricao TEXT,
			capa_url TEXT,
			pdf_url TEXT,
			disponivel BOOLEAN DEFAULT TRUE,
			media_nota NUMERIC(3,2) DEFAULT 0.0,
			total_avaliacoes INTEGER DEFAULT 0,
			paginas INTEGER DEFAULT 0,
			externo_id TEXT UNIQUE,
			fonte TEXT,
			status TEXT DEFAULT 'aprovado',
			curador_id INTEGER REFERENCES usuarios(id),
			search_vector tsvector
		);
	`)

	// Study & Social Tables
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS flashcards (
			id SERIAL PRIMARY KEY,
			usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
			material_id INTEGER REFERENCES materiais(id),
			pergunta TEXT NOT NULL,
			resposta TEXT NOT NULL,
			dificuldade INTEGER DEFAULT 0,
			proxima_revisao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS favoritos (
			usuario_id INTEGER REFERENCES usuarios(id),
			material_id INTEGER REFERENCES materiais(id),
			PRIMARY KEY (usuario_id, material_id)
		);
	`)
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS historico_leitura (
			id SERIAL PRIMARY KEY,
			usuario_id INTEGER REFERENCES usuarios(id),
			material_id INTEGER REFERENCES materiais(id),
			data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS interesses_usuario (
			usuario_id INTEGER REFERENCES usuarios(id),
			interesse TEXT,
			PRIMARY KEY (usuario_id, interesse)
		);
	`)

	// FTS Trigger & Index
	_, _ = db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_materiais_search_vector ON materiais USING gin(search_vector);
	`)
	_, _ = db.Exec(`
		CREATE OR REPLACE FUNCTION materiais_search_trigger() RETURNS trigger AS $$
		begin
		  new.search_vector :=
			setweight(to_tsvector('portuguese', coalesce(new.titulo,'')), 'A') ||
			setweight(to_tsvector('portuguese', coalesce(new.autor,'')), 'B') ||
			setweight(to_tsvector('portuguese', coalesce(new.descricao,'')), 'C');
		  return new;
		end
		$$ LANGUAGE plpgsql;
	`)
	_, _ = db.Exec(`
		DROP TRIGGER IF EXISTS tsvectorupdate ON materiais;
		CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
		ON materiais FOR EACH ROW EXECUTE FUNCTION materiais_search_trigger();
	`)

	// Additional Indexes for Performance
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_status ON materiais(status);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_categoria ON materiais(categoria);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_fonte ON materiais(fonte);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_ano ON materiais(ano_publicacao);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_externo_id ON materiais(externo_id) WHERE externo_id IS NOT NULL;`)

	geminiClient := ai.NewGeminiClient(cfg.GeminiAPIKey)

	// Inicia o worker de sincronização em segundo plano
	go startBackgroundSync(db)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Biblioteca Digital API"})
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		dbStatus := "up"
		if err := db.Ping(); err != nil {
			status = "degraded"
			dbStatus = "down"
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":   status,
			"database": dbStatus,
			"uptime":   time.Since(startTime).String(),
			"memory": map[string]interface{}{
				"alloc_mb": m.Alloc / 1024 / 1024,
				"total_mb": m.TotalAlloc / 1024 / 1024,
				"sys_mb":   m.Sys / 1024 / 1024,
				"num_gc":   m.NumGC,
			},
			"version": "1.3.0-expert",
		})
	})

	var c cache.Cache
	if cfg.RedisURL != "" {
		fmt.Printf("Usando Cache Redis em: %s\n", cfg.RedisURL)
		c = cache.NewRedisCache(cfg.RedisURL, cfg.RedisPassword)
	} else {
		fmt.Println("Usando Cache em Memória")
		c = cache.NewMemoryCache()
	}

	handler.RegisterUsuarioRoutes(mux, db)
	handler.RegisterMaterialRoutes(mux, db, geminiClient, c)

	handler.RegisterStatsRoutes(mux, db)
	handler.RegisterEstudoRoutes(mux, db, geminiClient)
	handler.RegisterAdminRoutes(mux, db)
	handler.RegisterAnotacaoRoutes(mux, db)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Apply Logger, Rate Limiter, and CORS middleware
	handlerWithSecurity := middleware.Security(mux)
	handlerWithLogger := middleware.Logger(handlerWithSecurity)
	handlerWithRateLimit := middleware.RateLimit(handlerWithLogger)
	handlerWithCORS := middleware.CORS(handlerWithRateLimit)

	logger.Info("Starting server", zap.String("port", cfg.Port))
	if err := http.ListenAndServe(":"+cfg.Port, handlerWithCORS); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

func startBackgroundSync(db *sql.DB) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	repo := &repository.MaterialPostgres{DB: db}
	mh := harvester.NewMultiSourceHarvester()

	// Executa uma vez no início
	logger.Info("Starting initial background synchronization")
	syncBooks(repo, mh)

	for range ticker.C {
		logger.Info("Executing periodic background synchronization")
		syncBooks(repo, mh)
	}
}

func syncBooks(repo *repository.MaterialPostgres, mh *harvester.MultiSourceHarvester) {
	if !syncMu.TryLock() {
		logger.Warn("Background synchronization already in progress, skipping this run")
		return
	}
	defer syncMu.Unlock()

	categories := []string{"TECNOLOGIA", "SAÚDE", "MATEMÁTICA", "CIÊNCIAS", "HISTÓRIA", "CONTABILIDADE"}
	for _, cat := range categories {
		mats, err := mh.Search(context.Background(), "", cat, "", 0, 0, 5)
		if err == nil {
			for i := range mats {
				_ = repo.Criar(context.Background(), &mats[i])
			}
		} else {
			logger.Error("Harvester search failed during sync", zap.String("category", cat), zap.Error(err))
		}
	}
	logger.Info("Background synchronization completed")
}
