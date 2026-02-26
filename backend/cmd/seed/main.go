package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type LivroMock struct {
	Nome          string `json:"nome"`
	AnoPublicacao int    `json:"ano_publicacao"`
	Autor         string `json:"autor"`
	Link          string `json:"link"`
}

func main() {
	if err := godotenv.Load(".env", "../../.env"); err != nil {
		log.Printf("Aviso: Erro ao carregar arquivo .env: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// tenta ler livros.json de varios locais relativos
	data, err := os.ReadFile("livros.json")
	if err != nil {
		data, err = os.ReadFile("../livros.json")
	}
	if err != nil {
		data, err = os.ReadFile("../../../livros.json")
	}

	if err != nil {
		log.Fatalf("Erro ao ler livros.json: %v", err)
	}

	var livros []LivroMock
	if err := json.Unmarshal(data, &livros); err != nil {
		log.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	fmt.Printf("Semeando %d livros...\n", len(livros))

	for _, l := range livros {
		isPDF := false
		if len(l.Link) > 4 && l.Link[len(l.Link)-4:] == ".pdf" {
			isPDF = true
		}

		var id int
		err := db.QueryRow("SELECT id FROM materiais WHERE titulo = $1", l.Nome).Scan(&id)

		if err == sql.ErrNoRows {
			fmt.Printf("Inserindo: %s\n", l.Nome)
			if isPDF {
				query := `INSERT INTO materiais (titulo, autor, categoria, ano_publicacao, pdf_url, disponivel)
						  VALUES ($1, $2, $3, $4, $5, $6)`
				_, err = db.Exec(query, l.Nome, l.Autor, "Tecnologia", l.AnoPublicacao, l.Link, true)
			} else {
				query := `INSERT INTO materiais (titulo, autor, categoria, ano_publicacao, capa_url, disponivel)
						  VALUES ($1, $2, $3, $4, $5, $6)`
				_, err = db.Exec(query, l.Nome, l.Autor, "Tecnologia", l.AnoPublicacao, l.Link, true)
			}
		} else if err == nil {
			fmt.Printf("Atualizando ID %d: %s\n", id, l.Nome)
			if isPDF {
				query := `UPDATE materiais SET pdf_url = $1, disponivel = $2 WHERE id = $3`
				_, err = db.Exec(query, l.Link, true, id)
			} else {
				query := `UPDATE materiais SET capa_url = $1, disponivel = $2 WHERE id = $3`
				_, err = db.Exec(query, l.Link, true, id)
			}
		} else {
			fmt.Printf("Erro ao buscar %s: %v\n", l.Nome, err)
		}

		if err != nil {
			log.Printf("Erro ao processar livro %s: %v", l.Nome, err)
		}
	}

	fmt.Println("Semeadura concluída com sucesso!")
}
