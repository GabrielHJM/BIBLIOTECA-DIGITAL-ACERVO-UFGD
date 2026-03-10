package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://postgres:admin@127.0.0.1:5432/BibliotecaDigital_BD?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Falha ao abrir conexão: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Falha ao pingar o banco (verifique se o Docker/Postgres está rodando): ", err)
	}

	queries := []string{
		"ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;",
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS data_criacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP;",
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;",
	}

	for _, q := range queries {
		fmt.Printf("Executando: %s\n", q)
		_, err := db.Exec(q)
		if err != nil {
			fmt.Printf("AVISO/ERRO: %v\n", err)
		} else {
			fmt.Println("Sucesso!")
		}
	}
	fmt.Println("Migração concluída.")
}
