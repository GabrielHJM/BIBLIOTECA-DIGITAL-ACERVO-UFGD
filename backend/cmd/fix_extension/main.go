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

	fmt.Println("Tentando habilitar a extensão unaccent...")
	_, err = db.Exec("CREATE EXTENSION IF NOT EXISTS unaccent")
	if err != nil {
		fmt.Printf("ERRO: %v\n", err)
		fmt.Println("\n--- DICA ---")
		fmt.Println("Se o erro for 'permission denied', pode ser necessário rodar como superusuário ou habilitar no Postgres.")
		fmt.Println("Se o erro for 'could not open extension control file', a extensão pode não estar instalada no OS.")
	} else {
		fmt.Println("SUCESSO: Extensão unaccent habilitada!")
	}
}
