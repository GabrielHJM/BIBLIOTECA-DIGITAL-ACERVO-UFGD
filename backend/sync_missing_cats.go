package main

import (
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgresql://postgres:postgres@localhost:5432/biblioteca_digital?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := &repository.MaterialPostgres{DB: db}
	multi := harvester.NewMultiSourceHarvester()

	categories := []string{"MATEMÁTICA", "HISTÓRIA"}

	for _, cat := range categories {
		fmt.Printf("Sincronizando categoria: %s...\n", cat)
		materials, err := multi.Search(context.Background(), "", cat, "", 0, 0, 10)
		if err != nil {
			fmt.Printf("Erro ao buscar %s: %v\n", cat, err)
			continue
		}

		fmt.Printf("Encontrados %d materiais para %s. Salvando...\n", len(materials), cat)
		savedCount := 0
		for i := range materials {
			materials[i].Categoria = cat // Garantir que está na categoria correta
			err := repo.Criar(context.Background(), &materials[i])
			if err == nil {
				savedCount++
			}
		}
		fmt.Printf("Salvos %d novos materiais para %s.\n", savedCount, cat)
	}

	fmt.Println("Sincronização concluída!")
	os.Exit(0)
}
