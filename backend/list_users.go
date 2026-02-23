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
	godotenv.Load(".env")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, nome, email FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("ID | Nome | Email")
	fmt.Println("---|------|------")
	for rows.Next() {
		var id int
		var nome, email string
		if err := rows.Scan(&id, &nome, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d | %s | %s\n", id, nome, email)
	}
}
