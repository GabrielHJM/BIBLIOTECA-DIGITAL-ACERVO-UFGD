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

	fmt.Println("Contagem de materiais por fonte:")
	rows, err := db.Query("SELECT fonte, COUNT(*) FROM materiais GROUP BY fonte")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var fonte string
		var count int
		if err := rows.Scan(&fonte, &count); err != nil {
			log.Fatal(err)
		}
		if fonte == "" {
			fonte = "Local/Desconhecida"
		}
		fmt.Printf("- %s: %d\n", fonte, count)
	}

	fmt.Println("\nCategorias distintas no banco:")
	rows, err = db.Query("SELECT DISTINCT categoria FROM materiais")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var cat string
			if err := rows.Scan(&cat); err == nil {
				fmt.Printf("- %s\n", cat)
			}
		}
	}

	fmt.Println("\nÚltimos 5 materiais inseridos:")
	rows, err = db.Query("SELECT titulo, fonte, status FROM materiais ORDER BY id DESC LIMIT 5")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var titulo, fonte, status string
		if err := rows.Scan(&titulo, &fonte, &status); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- [%s] %s (Status: %s)\n", fonte, titulo, status)
	}
}
