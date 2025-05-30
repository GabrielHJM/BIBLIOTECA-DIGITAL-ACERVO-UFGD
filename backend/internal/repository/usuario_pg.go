package repository

import (
	"database/sql"
	"biblioteca-digital-api/internal/domain/usuario"
)

type UsuarioPostgres struct {
	DB *sql.DB
}

func (r *UsuarioPostgres) Salvar(u *usuario.Usuario) error {
	// insert SQL
	return nil
}

func (r *UsuarioPostgres) BuscarPorEmail(email string) (*usuario.Usuario, error) {
	// select SQL
	return nil, nil
}

func (r *UsuarioPostgres) ListarInteresses(id int) ([]string, error) {
	// select SQL
	return nil, nil
}
