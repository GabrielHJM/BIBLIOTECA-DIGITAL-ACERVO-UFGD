package main

import (
	"biblioteca-digital-api/config"
	"log"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)
	defer db.Close()

	// Adiciona a coluna se ela não existir
	query := `ALTER TABLE usuarios ADD COLUMN IF NOT EXISTS meta_paginas_semana INTEGER DEFAULT 100;`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao adicionar coluna: %v\n", err)
	}

	log.Println("Migração concluída com sucesso: coluna 'meta_paginas_semana' adicionada/verificada.")
}
