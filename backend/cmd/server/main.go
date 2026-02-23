package main

import (
	"encoding/json"
	"log"
	"net/http"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"
	"biblioteca-digital-api/internal/handler/middleware"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "message": "Biblioteca Digital API"})
	})

	handler.RegisterUsuarioRoutes(mux, db)
	handler.RegisterMaterialRoutes(mux, db)

	// Apply Logger and CORS middleware
	handlerWithLogger := middleware.Logger(mux)
	handlerWithCORS := middleware.CORS(handlerWithLogger)

	log.Printf("Servidor rodando na porta %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handlerWithCORS); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
