package main

import (
	"biblioteca-digital-api/config"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)
	defer db.Close()

	rows, err := db.QueryContext(context.Background(), "SELECT nome, email FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Usuários cadastrados:")
	for rows.Next() {
		var nome, email string
		if err := rows.Scan(&nome, &email); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("- %s (%s)\n", nome, email)
	}
}
