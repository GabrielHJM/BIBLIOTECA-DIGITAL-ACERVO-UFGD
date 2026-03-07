package repository

import (
	"biblioteca-digital-api/internal/domain/anotacao"
	"context"
	"database/sql"
	"time"
)

type AnotacaoRepositoryPG struct {
	db *sql.DB
}

func NewAnotacaoRepositoryPG(db *sql.DB) *AnotacaoRepositoryPG {
	return &AnotacaoRepositoryPG{db: db}
}

func (r *AnotacaoRepositoryPG) Create(ctx context.Context, req anotacao.Anotacao) (int, error) {
	query := `
		INSERT INTO anotacoes (usuario_id, material_id, titulo, conteudo, cor)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	var id int
	err := r.db.QueryRowContext(ctx, query, req.UsuarioID, req.MaterialID, req.Titulo, req.Conteudo, req.Cor).Scan(&id)
	return id, err
}

func (r *AnotacaoRepositoryPG) GetByID(ctx context.Context, id int) (anotacao.Anotacao, error) {
	query := `
		SELECT id, usuario_id, material_id, titulo, conteudo, cor, data_criacao, data_atualizacao
		FROM anotacoes
		WHERE id = $1
	`
	var a anotacao.Anotacao
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.UsuarioID, &a.MaterialID, &a.Titulo, &a.Conteudo, &a.Cor, &a.DataCriacao, &a.DataAtualizacao,
	)
	return a, err
}

func (r *AnotacaoRepositoryPG) ListByUsuario(ctx context.Context, usuarioID int) ([]anotacao.Anotacao, error) {
	query := `
		SELECT id, usuario_id, material_id, titulo, conteudo, cor, data_criacao, data_atualizacao
		FROM anotacoes
		WHERE usuario_id = $1
		ORDER BY data_atualizacao DESC
	`
	rows, err := r.db.QueryContext(ctx, query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []anotacao.Anotacao
	for rows.Next() {
		var a anotacao.Anotacao
		if err := rows.Scan(&a.ID, &a.UsuarioID, &a.MaterialID, &a.Titulo, &a.Conteudo, &a.Cor, &a.DataCriacao, &a.DataAtualizacao); err != nil {
			return nil, err
		}
		lista = append(lista, a)
	}
	return lista, nil
}

func (r *AnotacaoRepositoryPG) Update(ctx context.Context, req anotacao.Anotacao) error {
	query := `
		UPDATE anotacoes
		SET titulo = $1, conteudo = $2, cor = $3, data_atualizacao = $4
		WHERE id = $5 AND usuario_id = $6
	`
	_, err := r.db.ExecContext(ctx, query, req.Titulo, req.Conteudo, req.Cor, time.Now(), req.ID, req.UsuarioID)
	return err
}

func (r *AnotacaoRepositoryPG) Delete(ctx context.Context, id int, usuarioID int) error {
	query := `
		DELETE FROM anotacoes
		WHERE id = $1 AND usuario_id = $2
	`
	_, err := r.db.ExecContext(ctx, query, id, usuarioID)
	return err
}
