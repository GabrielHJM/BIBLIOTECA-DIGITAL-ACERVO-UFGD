package repository

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"context"
	"database/sql"
	"errors"
	"strings"
)

type UsuarioPostgres struct {
	DB *sql.DB
}

func NewUsuarioPG(db *sql.DB) *UsuarioPostgres {
	return &UsuarioPostgres{DB: db}
}

func (r *UsuarioPostgres) Salvar(ctx context.Context, u *usuario.Usuario) error {
	query := "INSERT INTO usuarios (nome, email, senha, tipo, foto_url, cpf, data_nascimento, username) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	err := r.DB.QueryRowContext(ctx, query, u.Nome, u.Email, u.Senha, u.Tipo, u.FotoURL, u.Cpf, u.DataNascimento, u.Username).Scan(&u.ID)
	if err != nil {
		// Verificar se é erro de violação de unicidade (email ou username duplicado)
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "username") {
				return errors.New("este nome de usuário já está sendo usado")
			}
			return errors.New("este email já está cadastrado")
		}
		return err
	}
	return nil
}

func (r *UsuarioPostgres) BuscarPorEmail(ctx context.Context, email string) (*usuario.Usuario, error) {
	query := "SELECT id, nome, email, senha, COALESCE(tipo, 1), COALESCE(foto_url, ''), cpf, data_nascimento, COALESCE(username, ''), deleted_at FROM usuarios WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL"
	u := &usuario.Usuario{}
	err := r.DB.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Nome, &u.Email, &u.Senha, &u.Tipo, &u.FotoURL, &u.Cpf, &u.DataNascimento, &u.Username, &u.DeletedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UsuarioPostgres) BuscarPorID(ctx context.Context, id int) (*usuario.Usuario, error) {
	query := "SELECT id, nome, email, senha, COALESCE(tipo, 1), COALESCE(foto_url, ''), cpf, data_nascimento, COALESCE(username, ''), deleted_at FROM usuarios WHERE id = $1 AND deleted_at IS NULL"
	u := &usuario.Usuario{}
	err := r.DB.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Nome, &u.Email, &u.Senha, &u.Tipo, &u.FotoURL, &u.Cpf, &u.DataNascimento, &u.Username, &u.DeletedAt)
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
	query := "UPDATE usuarios SET senha = $1 WHERE LOWER(email) = LOWER($2)"
	_, err := r.DB.ExecContext(ctx, query, novaSenha, email)
	return err
}

func (r *UsuarioPostgres) Atualizar(ctx context.Context, u *usuario.Usuario) error {
	query := "UPDATE usuarios SET nome = $1, email = $2, foto_url = $3, cpf = $4, data_nascimento = $5, username = $6, senha = $7 WHERE id = $8 AND deleted_at IS NULL"
	_, err := r.DB.ExecContext(ctx, query, u.Nome, u.Email, u.FotoURL, u.Cpf, u.DataNascimento, u.Username, u.Senha, u.ID)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "23505") {
			if strings.Contains(err.Error(), "username") {
				return errors.New("este nome de usuário já está sendo usado")
			}
			return errors.New("este email já está sendo usado")
		}
		return err
	}
	return err
}

func (r *UsuarioPostgres) AtualizarMeta(ctx context.Context, id int, meta int) error {
	query := "UPDATE usuarios SET meta_paginas_semana = $1 WHERE id = $2 AND deleted_at IS NULL"
	_, err := r.DB.ExecContext(ctx, query, meta, id)
	return err
}

func (r *UsuarioPostgres) Deletar(ctx context.Context, id int) error {
	query := "UPDATE usuarios SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1"
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}
