package main

import (
	"biblioteca-digital-api/config"
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

	fmt.Println("Listando materiais recentes:")
	// Using custom query to get the latest 5 materials
	query := "SELECT id, titulo, fonte, capa_url, dificuldade, xp, pdf_url FROM materiais ORDER BY id DESC LIMIT 5"
	rows, err := db.QueryContext(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, dificuldade, xp int
		var titulo, fonte, capaURL, pdfURL string
		if err := rows.Scan(&id, &titulo, &fonte, &capaURL, &dificuldade, &xp, &pdfURL); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- Titulo: %s\n  Fonte: %s\n  Capa: %s\n  PDF: %s\n  Dificuldade: %d, XP: %d\n\n", titulo, fonte, capaURL, pdfURL, dificuldade, xp)
	}
}
