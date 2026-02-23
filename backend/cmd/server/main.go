package main

import (
	"log"
	"net/http"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente do sistema")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)

	mux := http.NewServeMux()
	handler.RegisterUsuarioRoutes(mux, db)

	log.Printf("Servidor rodando na porta %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Erro ao iniciar o servidor: %v", err)
	}
}
