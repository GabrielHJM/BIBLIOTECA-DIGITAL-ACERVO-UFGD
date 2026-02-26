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
	godotenv.Load(".env", "../../.env")
	dbUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, titulo, pdf_url, disponivel FROM materiais WHERE id = 1")

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var titulo string
		rows.Scan(&id, &titulo)
		fmt.Printf("[%d] -> %q\n", id, titulo)
	}
}
