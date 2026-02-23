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

	fmt.Println("=== Schema for 'usuarios' ===")
	rows, err := db.Query("SELECT column_name, data_type, is_nullable FROM information_schema.columns WHERE table_name = 'usuarios'")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name, dtype, nullable string
		rows.Scan(&name, &dtype, &nullable)
		fmt.Printf("%s (%s) - Nullable: %s\n", name, dtype, nullable)
	}
	rows.Close()

	fmt.Println("\n=== All Users ===")
	rows, err = db.Query("SELECT id, nome, email, senha FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var nome, email, senha string
		if err := rows.Scan(&id, &nome, &email, &senha); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d | Nome: %s | Email: %s | Hash: %s\n", id, nome, email, senha[:10]+"...")
	}
}
