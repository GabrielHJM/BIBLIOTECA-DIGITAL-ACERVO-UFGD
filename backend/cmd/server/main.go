package main

import (
	"log"
	"net/http"

	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/handler"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)
	handler := handler.NewHandler(db)

	http.HandleFunc("/api/usuarios", handler.UsuarioHandler)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}