package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"
	"biblioteca-digital-api/internal/handler/middleware"
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"

	"runtime"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var startTime = time.Now()

func main() {
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)
	defer db.Close()

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

	handler.RegisterUsuarioRoutes(mux, db)
	handler.RegisterMaterialRoutes(mux, db)

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
