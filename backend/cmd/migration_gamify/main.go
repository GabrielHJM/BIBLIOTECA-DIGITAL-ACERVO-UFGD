package main

import (
	"biblioteca-digital-api/config"
	"database/sql"
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
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("Aplicando colunas de gamificação...")

	queries := []string{
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS dificuldade INTEGER DEFAULT 1",
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS xp INTEGER DEFAULT 10",
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS relevancia INTEGER DEFAULT 0",
	}

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			fmt.Printf("ERRO ao executar [%s]: %v\n", q, err)
		} else {
			fmt.Printf("SUCESSO: %s\n", q)
		}
	}

	fmt.Println("Migração concluída.")
}
