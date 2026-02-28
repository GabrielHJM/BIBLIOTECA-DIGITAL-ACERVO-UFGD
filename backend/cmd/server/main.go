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
	"biblioteca-digital-api/internal/usecase/social"
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
	defer db.Close()

	// Gamification Migration
	_, _ = db.Exec(`ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS meta_paginas_semana INTEGER DEFAULT 100;`)

	// Study Tools Migrations
	_, _ = db.Exec(`ALTER TABLE materiais ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'aprovado';`)
	_, _ = db.Exec(`ALTER TABLE materiais ADD COLUMN IF NOT EXISTS curador_id INTEGER REFERENCES usuarios(id);`)

	// Optimization Indexes
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_status ON materiais(status);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_categoria ON materiais(categoria);`)
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_externo_id ON materiais(externo_id) WHERE externo_id IS NOT NULL;`)

	_, _ = db.Exec(`
		CREATE TABLE IF NOT EXISTS anotacoes (
			id SERIAL PRIMARY KEY,
			usuario_id INTEGER NOT NULL REFERENCES usuarios(id),
			material_id INTEGER NOT NULL REFERENCES materiais(id),
			conteudo TEXT NOT NULL,
			pagina INTEGER,
			cor TEXT DEFAULT 'yellow',
			data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
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

	socialRepo := repository.NewSocialPG(db)
	socialUC := social.NewSocialUseCase(socialRepo, socialRepo, socialRepo, socialRepo)
	handler.RegisterSocialRoutes(mux, socialUC, socialRepo, socialRepo)

	handler.RegisterStatsRoutes(mux, db)
	handler.RegisterEstudoRoutes(mux, db, geminiClient)
	handler.RegisterAdminRoutes(mux, db)

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// Apply Logger and CORS middleware
	handlerWithLogger := middleware.Logger(mux)
	handlerWithCORS := middleware.CORS(handlerWithLogger)

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
