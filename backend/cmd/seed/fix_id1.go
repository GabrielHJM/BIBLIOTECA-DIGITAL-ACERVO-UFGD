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

	_, err = db.Exec("UPDATE materiais SET pdf_url = $1, disponivel = true WHERE id = 1", "https://www.ifsc.usp.br/~igor/teaching/fiscomp1/ood-java.pdf")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("ID 1 updated manually.")
}
