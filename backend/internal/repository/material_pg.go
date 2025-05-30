package repository

import (
	"database/sql"
	"biblioteca-digital-api/internal/domain/material"
)

type MaterialPostgres struct {
	DB *sql.DB
}

func (r *MaterialPostgres) Listar() ([]material.Material, error) {
	// select SQL
	return nil, nil
}

func (r *MaterialPostgres) Favoritar(id int) error {
	// update SQL
	return nil
}
