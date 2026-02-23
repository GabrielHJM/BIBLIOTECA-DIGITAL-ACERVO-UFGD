package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const schema = `
-- Tabela de Usuários
CREATE TABLE IF NOT EXISTS usuarios (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    senha TEXT NOT NULL,
    tipo TEXT NOT NULL DEFAULT 'estudante'
);

-- Tabela de Materiais (Livros/Artigos)
CREATE TABLE IF NOT EXISTS materiais (
    id SERIAL PRIMARY KEY,
    titulo TEXT NOT NULL,
    autor TEXT NOT NULL,
    isbn TEXT,
    categoria TEXT,
    ano_publicacao INTEGER,
    descricao TEXT,
    capa_url TEXT,
    pdf_url TEXT,
    disponivel BOOLEAN DEFAULT TRUE,
    media_nota DECIMAL(3,2) DEFAULT 0,
    total_avaliacoes INTEGER DEFAULT 0
);

-- Tabela de Interesses do Usuário
CREATE TABLE IF NOT EXISTS interesses_usuario (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    interesse TEXT NOT NULL
);

-- Tabela de Favoritos
CREATE TABLE IF NOT EXISTS favoritos (
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    PRIMARY KEY (usuario_id, material_id)
);

-- Tabela de Avaliações
CREATE TABLE IF NOT EXISTS avaliacoes (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    nota INTEGER CHECK (nota >= 1 AND nota <= 5),
    comentario TEXT,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Empréstimos
CREATE TABLE IF NOT EXISTS emprestimos (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data_emprestimo TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_devolucao TIMESTAMP,
    status TEXT DEFAULT 'ativo'
);

-- Tabela de Histórico de Leitura
CREATE TABLE IF NOT EXISTS historico_leitura (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Erro ao carregar arquivo .env: %v", err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL não configurada no .env")
	}

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer db.Close()

	fmt.Println("Iniciando criação do schema...")
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Erro ao criar tabelas: %v", err)
	}

	fmt.Println("Schema criado com sucesso! O banco de dados está pronto para uso.")
}
