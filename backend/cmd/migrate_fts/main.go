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
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Aviso:", err)
	}
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL não configurada no .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Erro ao conectar: %v", err)
	}
	defer db.Close()

	queries := []string{
		// 1. Adicionar coluna tsvector
		"ALTER TABLE materiais ADD COLUMN IF NOT EXISTS search_vector tsvector;",

		// 2. Criar função para atualizar o search_vector
		`CREATE OR REPLACE FUNCTION materiais_search_trigger() RETURNS trigger AS $$
		begin
		  new.search_vector :=
		    setweight(to_tsvector('portuguese', coalesce(new.titulo,'')), 'A') ||
		    setweight(to_tsvector('portuguese', coalesce(new.autor,'')), 'B') ||
		    setweight(to_tsvector('portuguese', coalesce(new.descricao,'')), 'C');
		  return new;
		end
		$$ LANGUAGE plpgsql;`,

		// 3. Criar trigger se não existir
		`DO $$
		BEGIN
		    IF NOT EXISTS (SELECT 1 FROM pg_trigger WHERE tgname = 'trg_materiais_search') THEN
		        CREATE TRIGGER trg_materiais_search BEFORE INSERT OR UPDATE
		        ON materiais FOR EACH ROW EXECUTE FUNCTION materiais_search_trigger();
		    END IF;
		END
		$$;`,

		// 4. Criar índice GIN para busca rápida
		"CREATE INDEX IF NOT EXISTS idx_materiais_search ON materiais USING GIN(search_vector);",

		// 5. Popular dados existentes
		"UPDATE materiais SET search_vector = setweight(to_tsvector('portuguese', coalesce(titulo,'')), 'A') || setweight(to_tsvector('portuguese', coalesce(autor,'')), 'B') || setweight(to_tsvector('portuguese', coalesce(descricao,'')), 'C') WHERE search_vector IS NULL;",
	}

	for _, q := range queries {
		_, err = db.Exec(q)
		if err != nil {
			log.Printf("Erro na query: %v\n", err)
		} else {
			fmt.Printf("Sucesso na migração de FTS.\n")
		}
	}
	fmt.Println("Migração FTS concluída.")
}
