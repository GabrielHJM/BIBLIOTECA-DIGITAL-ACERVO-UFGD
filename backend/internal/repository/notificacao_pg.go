package repository

import (
	"biblioteca-digital-api/internal/domain/notificacao"
	"context"
	"database/sql"
)

type NotificacaoPostgres struct {
	DB *sql.DB
}

func NewNotificacaoPostgres(db *sql.DB) *NotificacaoPostgres {
	return &NotificacaoPostgres{DB: db}
}

func (r *NotificacaoPostgres) Criar(ctx context.Context, n *notificacao.Notificacao) error {
	query := `INSERT INTO notificacoes (usuario_id, titulo, mensagem, tipo)
	          VALUES ($1, $2, $3, $4) RETURNING id, data_criacao`
	return r.DB.QueryRowContext(ctx, query, n.UsuarioID, n.Titulo, n.Mensagem, n.Tipo).
		Scan(&n.ID, &n.DataCriacao)
}

func (r *NotificacaoPostgres) ListarPorUsuario(ctx context.Context, usuarioID int) ([]notificacao.Notificacao, error) {
	query := `SELECT id, usuario_id, titulo, mensagem, tipo, lida, data_criacao
	          FROM notificacoes WHERE usuario_id = $1 ORDER BY data_criacao DESC`
	rows, err := r.DB.QueryContext(ctx, query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ns []notificacao.Notificacao
	for rows.Next() {
		var n notificacao.Notificacao
		if err := rows.Scan(&n.ID, &n.UsuarioID, &n.Titulo, &n.Mensagem, &n.Tipo, &n.Lida, &n.DataCriacao); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

func (r *NotificacaoPostgres) MarcarComoLida(ctx context.Context, id int) error {
	query := `UPDATE notificacoes SET lida = TRUE WHERE id = $1`
	_, err := r.DB.ExecContext(ctx, query, id)
	return err
}

func (r *NotificacaoPostgres) LimparPorUsuario(ctx context.Context, usuarioID int) error {
	query := `DELETE FROM notificacoes WHERE usuario_id = $1`
	_, err := r.DB.ExecContext(ctx, query, usuarioID)
	return err
}
