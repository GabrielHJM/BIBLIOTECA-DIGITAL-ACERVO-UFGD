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
	"strings"
	"time"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"
	"biblioteca-digital-api/internal/handler/middleware"
	"biblioteca-digital-api/internal/harvester"
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

	// Configurar Timezone da Sessão
	if _, err := db.Exec("SET TIMEZONE TO 'America/Sao_Paulo';"); err != nil {
		logger.Error("Failed to set session timezone", zap.Error(err))
	}

	// Core Migrations
	if _, err := db.Exec("CREATE EXTENSION IF NOT EXISTS unaccent;"); err != nil {
		logger.Error("Failed to create unaccent extension", zap.Error(err))
	}

	migrations := []struct {
		name string
		sql  string
	}{
		{"usuarios", `CREATE TABLE IF NOT EXISTS usuarios (
			id SERIAL PRIMARY KEY, nome TEXT NOT NULL, email TEXT UNIQUE NOT NULL, senha TEXT NOT NULL,
			tipo INTEGER DEFAULT 1, foto_url TEXT, cpf TEXT, data_nascimento DATE, username TEXT UNIQUE,
			meta_paginas_semana INTEGER DEFAULT 100, data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP
		);`},
		{"materiais", `CREATE TABLE IF NOT EXISTS materiais (
			id SERIAL PRIMARY KEY, titulo TEXT NOT NULL, autor TEXT NOT NULL, isbn TEXT, categoria TEXT NOT NULL,
			ano_publicacao INTEGER, descricao TEXT, capa_url TEXT, pdf_url TEXT, disponivel BOOLEAN DEFAULT TRUE,
			media_nota NUMERIC(3,2) DEFAULT 0.0, total_avaliacoes INTEGER DEFAULT 0, paginas INTEGER DEFAULT 0,
			externo_id TEXT UNIQUE, fonte TEXT, status TEXT DEFAULT 'aprovado', curador_id INTEGER REFERENCES usuarios(id),
			dificuldade INTEGER DEFAULT 1, xp INTEGER DEFAULT 10, relevancia INTEGER DEFAULT 0, search_vector tsvector,
			data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP, deleted_at TIMESTAMP
		);`},
		{"flashcards", `CREATE TABLE IF NOT EXISTS flashcards (
			id SERIAL PRIMARY KEY, usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
			material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL, pergunta TEXT NOT NULL,
			resposta TEXT NOT NULL, dificuldade INTEGER DEFAULT 0, proxima_revisao TIMESTAMP DEFAULT CURRENT_TIMESTAMP, data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"favoritos", `CREATE TABLE IF NOT EXISTS favoritos (
			usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE, PRIMARY KEY (usuario_id, material_id)
		);`},
		{"historico", `CREATE TABLE IF NOT EXISTS historico_leitura (
			id SERIAL PRIMARY KEY, usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"interesses", `CREATE TABLE IF NOT EXISTS interesses_usuario (
			usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, interesse TEXT NOT NULL, PRIMARY KEY (usuario_id, interesse)
		);`},
		{"notificacoes", `CREATE TABLE IF NOT EXISTS notificacoes (
			id SERIAL PRIMARY KEY, usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE,
			titulo TEXT NOT NULL, mensagem TEXT NOT NULL, tipo TEXT DEFAULT 'info', lida BOOLEAN DEFAULT FALSE, data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"avaliacoes", `CREATE TABLE IF NOT EXISTS avaliacoes (
			id SERIAL PRIMARY KEY, usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
			nota INTEGER CHECK (nota >= 1 AND nota <= 5), comentario TEXT, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"emprestimos", `CREATE TABLE IF NOT EXISTS emprestimos (
			id SERIAL PRIMARY KEY, usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
			data_emprestimo TIMESTAMP DEFAULT CURRENT_TIMESTAMP, data_devolucao TIMESTAMP, status TEXT DEFAULT 'ativo'
		);`},
		{"curtidas", `CREATE TABLE IF NOT EXISTS curtidas (
			usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
			data TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (usuario_id, material_id)
		);`},
		{"comentarios", `CREATE TABLE IF NOT EXISTS comentarios (
			id SERIAL PRIMARY KEY, usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
			texto TEXT NOT NULL, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"amizades", `CREATE TABLE IF NOT EXISTS amizades (
			usuario_id1 INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, usuario_id2 INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
			status TEXT DEFAULT 'pendente', data_solicitacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (usuario_id1, usuario_id2), CHECK (usuario_id1 != usuario_id2)
		);`},
		{"anotacoes", `CREATE TABLE IF NOT EXISTS anotacoes (
			id SERIAL PRIMARY KEY, usuario_id INTEGER NOT NULL REFERENCES usuarios(id) ON DELETE CASCADE, material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL,
			titulo TEXT, conteudo TEXT NOT NULL, cor TEXT DEFAULT '#FFFFFF', data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP, data_atualizacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"mensagens", `CREATE TABLE IF NOT EXISTS mensagens (
			id SERIAL PRIMARY KEY, remetente_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE, destinatario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
			material_id INTEGER REFERENCES materiais(id) ON DELETE SET NULL, conteudo TEXT, lida BOOLEAN DEFAULT FALSE, data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`},
		{"materiais_data_criacao", `ALTER TABLE materiais ADD COLUMN IF NOT EXISTS data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`},
		{"materiais_deleted_at", `ALTER TABLE materiais ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;`},
		{"usuarios_deleted_at", `ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;`},
		{"usuarios_cpf", `ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS cpf TEXT;`},
		{"usuarios_data_nascimento", `ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS data_nascimento DATE;`},
		{"usuarios_username", `ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS username TEXT UNIQUE;`},
	}

	for _, m := range migrations {
		if _, err := db.Exec(m.sql); err != nil {
			logger.Error("Migration failed", zap.String("table", m.name), zap.Error(err))
		}
	}

	// FTS Trigger & Index
	if _, err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_materiais_search_vector ON materiais USING gin(search_vector);`); err != nil {
		logger.Error("Failed to create FTS index", zap.Error(err))
	}
	if _, err := db.Exec(`
		CREATE OR REPLACE FUNCTION materiais_search_trigger() RETURNS trigger AS $$
		begin
		  new.search_vector :=
			setweight(to_tsvector('portuguese', unaccent(coalesce(new.titulo,''))), 'A') ||
			setweight(to_tsvector('portuguese', unaccent(coalesce(new.autor,''))), 'B') ||
			setweight(to_tsvector('portuguese', unaccent(coalesce(new.descricao,''))), 'C');
		  return new;
		end
		$$ LANGUAGE plpgsql;
	`); err != nil {
		logger.Error("Failed to create FTS trigger function", zap.Error(err))
	}
	if _, err := db.Exec(`
		DROP TRIGGER IF EXISTS tsvectorupdate ON materiais;
		CREATE TRIGGER tsvectorupdate BEFORE INSERT OR UPDATE
		ON materiais FOR EACH ROW EXECUTE FUNCTION materiais_search_trigger();
	`); err != nil {
		logger.Error("Failed to create FTS trigger", zap.Error(err))
	}

	// Additional Indexes for Performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_materiais_status ON materiais(status);",
		"CREATE INDEX IF NOT EXISTS idx_materiais_categoria ON materiais(categoria);",
		"CREATE INDEX IF NOT EXISTS idx_materiais_fonte ON materiais(fonte);",
		"CREATE INDEX IF NOT EXISTS idx_materiais_ano ON materiais(ano_publicacao);",
		"CREATE INDEX IF NOT EXISTS idx_materiais_externo_id ON materiais(externo_id) WHERE externo_id IS NOT NULL;",
	}
	for _, idx := range indexes {
		if _, err := db.Exec(idx); err != nil {
			logger.Error("Index creation failed", zap.String("sql", idx), zap.Error(err))
		}
	}
	_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_historico_usuario ON historico_leitura(usuario_id);`)

	// Inicia o worker de sincronização em segundo plano
	go startBackgroundSync(db)

	mux := http.NewServeMux()
	// Serve static files from the "dist" directory
	fs := http.FileServer(http.Dir("dist"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If the request is for an API or Swagger, let the other handlers deal with it
		// (ServeMux matches the most specific pattern, so this is the fallback)
		
		// Check if file exists in dist
		path := r.URL.Path
		if path == "/" {
			fs.ServeHTTP(w, r)
			return
		}

		// Avoid serving "dist" files for /api routes if they somehow fall through
		if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/swagger") || strings.HasPrefix(path, "/health") {
			http.NotFound(w, r)
			return
		}

		// Check if file exists in dist
		relativePath := strings.TrimPrefix(path, "/")
		f, err := http.Dir("dist").Open(relativePath)
		if err == nil {
			f.Close()
			fs.ServeHTTP(w, r)
			return
		}

		// File NOT found in dist.
		// If the request is for a static asset (has an extension), return 404 to avoid serving HTML as JS.
		// Otherwise (SPA route), serve index.html for client-side routing.
		if strings.Contains(path, ".") && !strings.HasSuffix(path, ".html") {
			logger.Warn("Static file not found", zap.String("path", path))
			http.NotFound(w, r)
			return
		}

		// Fallback to index.html for SPA routing
		http.ServeFile(w, r, "dist/index.html")
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

	// Create a sub-mux for API routes
	apiMux := http.NewServeMux()
	handler.RegisterUsuarioRoutes(apiMux, db)
	handler.RegisterMaterialRoutes(apiMux, db, c)
	handler.RegisterStatsRoutes(apiMux, db)
	handler.RegisterEstudoRoutes(apiMux, db)
	handler.RegisterAdminRoutes(apiMux, db)
	handler.RegisterAnotacaoRoutes(apiMux, db)
	handler.RegisterNotificacaoRoutes(apiMux, db)

	// Mount the API sub-mux at /api
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

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

	// High-yield categories for initial population
	categories := []string{
		"TECNOLOGIA", "SAÚDE", "DIREITO", "CIÊNCIAS", "MATEMÁTICA", "EDUCAÇÃO", 
		"LITERATURA BRASILEIRA", "HISTÓRIA DO BRASIL", "CONTABILIDADE",
	}
	
	logger.Info("Sync: Starting high-volume harvest", zap.Int("categories_count", len(categories)))

	for _, cat := range categories {
		// Fetch more books per category for a richer initial experience
		mats, err := mh.Search(context.Background(), "", cat, "", 0, 0, 15)
		if err == nil {
			count := 0
			for i := range mats {
				if err := repo.Criar(context.Background(), &mats[i]); err != nil {
					if !strings.Contains(err.Error(), "já existe") {
						logger.Debug("Sync: Failed to save material", zap.String("title", mats[i].Titulo), zap.Error(err))
					}
				} else {
					count++
				}
			}
			logger.Info("Sync: Category populated", zap.String("category", cat), zap.Int("new_books", count))
		} else {
			logger.Error("Harvester search failed during sync", zap.String("category", cat), zap.Error(err))
		}
	}
	logger.Info("Background synchronization completed successfully")
}
