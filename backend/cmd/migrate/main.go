package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Aviso:", err)
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL não configurada no .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer db.Close()

	queries := []string{
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS externo_id TEXT DEFAULT '';",
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS fonte TEXT DEFAULT '';",
		"DROP INDEX IF EXISTS unique_externo_id_idx;",
		"CREATE UNIQUE INDEX IF NOT EXISTS unique_externo_id_idx ON materiais (externo_id) WHERE externo_id != '';",
	}

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			log.Printf("Erro na query '%s': %v\n", q, err)
		} else {
			fmt.Printf("Sucesso: %s\n", q)
		}
	}
	fmt.Println("Migração concluída.")
}
