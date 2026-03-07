package main

import (
	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/repository"
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)
	repo := &repository.MaterialPostgres{DB: db}

	fmt.Println("Testando Pesquisar(categoria='TECNOLOGIA'):")
	materials, err := repo.Pesquisar(context.Background(), "", "TECNOLOGIA", "", 0, 0, nil, 10, 0, "")
	if err != nil {
		fmt.Printf("ERRO DETECTADO: %v\n", err)
	} else {
		fmt.Printf("SUCESSO: %d materiais encontrados\n", len(materials))
	}
}
