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
-- Drop existing tables for a clean start
DROP TABLE IF EXISTS emprestimos CASCADE;
DROP TABLE IF EXISTS historico_leitura CASCADE;
DROP TABLE IF EXISTS avaliacoes CASCADE;
DROP TABLE IF EXISTS favoritos CASCADE;
DROP TABLE IF EXISTS interesses_usuario CASCADE;
DROP TABLE IF EXISTS materiais CASCADE;
DROP TABLE IF EXISTS usuarios CASCADE;

-- Tabela de Usuários
CREATE TABLE usuarios (
    id SERIAL PRIMARY KEY,
    nome TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    senha TEXT NOT NULL,
    tipo TEXT NOT NULL DEFAULT 'estudante',
    foto_url TEXT
);

-- Tabela de Materiais (Livros/Artigos)
CREATE TABLE materiais (
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
CREATE TABLE interesses_usuario (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    interesse TEXT NOT NULL
);

-- Tabela de Favoritos
CREATE TABLE favoritos (
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    PRIMARY KEY (usuario_id, material_id)
);

-- Tabela de Avaliações
CREATE TABLE avaliacoes (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    nota INTEGER CHECK (nota >= 1 AND nota <= 5),
    comentario TEXT,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabela de Empréstimos
CREATE TABLE emprestimos (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data_emprestimo TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_devolucao TIMESTAMP,
    status TEXT DEFAULT 'ativo'
);

-- Tabela de Histórico de Leitura
CREATE TABLE historico_leitura (
    id SERIAL PRIMARY KEY,
    usuario_id INTEGER REFERENCES usuarios(id) ON DELETE CASCADE,
    material_id INTEGER REFERENCES materiais(id) ON DELETE CASCADE,
    data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func main() {
	// tenta carregar o .env da raiz do backend ou de dois niveis acima
	if err := godotenv.Load(".env", "../../.env"); err != nil {
		log.Printf("Aviso: Erro ao carregar arquivo .env: %v. Continuando com variáveis de ambiente.", err)
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
