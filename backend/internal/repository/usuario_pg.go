package repository

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"context"
	"database/sql"
)

type UsuarioPostgres struct {
	DB *sql.DB
}

func NewUsuarioPG(db *sql.DB) *UsuarioPostgres {
	return &UsuarioPostgres{DB: db}
}

func (r *UsuarioPostgres) Salvar(ctx context.Context, u *usuario.Usuario) error {
	query := "INSERT INTO usuarios (nome, email, senha, tipo, foto_url) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return r.DB.QueryRowContext(ctx, query, u.Nome, u.Email, u.Senha, u.Tipo, u.FotoURL).Scan(&u.ID)
}

func (r *UsuarioPostgres) BuscarPorEmail(ctx context.Context, email string) (*usuario.Usuario, error) {
	query := "SELECT id, nome, email, senha, COALESCE(tipo, 1), COALESCE(foto_url, '') FROM usuarios WHERE email = $1"
	u := &usuario.Usuario{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Nome, &u.Email, &u.Senha, &u.Tipo, &u.FotoURL)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UsuarioPostgres) ListarInteresses(ctx context.Context, id int) ([]string, error) {
	query := "SELECT interesse FROM interesses_usuario WHERE usuario_id = $1"
	rows, err := r.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var interesses []string
	for rows.Next() {
		var interesse string
		if err := rows.Scan(&interesse); err != nil {
			return nil, err
		}
		interesses = append(interesses, interesse)
	}
	return interesses, nil
}

func (r *UsuarioPostgres) AtualizarSenha(ctx context.Context, email string, novaSenha string) error {
	query := "UPDATE usuarios SET senha = $1 WHERE email = $2"
	_, err := r.DB.ExecContext(ctx, query, novaSenha, email)
	return err
}

func (r *UsuarioPostgres) Atualizar(ctx context.Context, u *usuario.Usuario) error {
	query := "UPDATE usuarios SET nome = $1, email = $2, foto_url = $3 WHERE id = $4"
	_, err := r.DB.ExecContext(ctx, query, u.Nome, u.Email, u.FotoURL, u.ID)
	return err
}

func (r *UsuarioPostgres) Deletar(ctx context.Context, id int) error {
	query := "DELETE FROM usuarios WHERE id = $1"
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
