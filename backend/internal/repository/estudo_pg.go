package repository

import (
	"biblioteca-digital-api/internal/domain/estudo"
	"context"
	"database/sql"
	"time"
)

type EstudoPostgres struct {
	DB *sql.DB
}

func NewEstudoPostgres(db *sql.DB) *EstudoPostgres {
	return &EstudoPostgres{DB: db}
}

func (r *EstudoPostgres) CriarFlashcard(ctx context.Context, f *estudo.Flashcard) error {
	query := `INSERT INTO flashcards (usuario_id, material_id, pergunta, resposta)
	          VALUES ($1, $2, $3, $4) RETURNING id, proxima_revisao, data_criacao`
	return r.DB.QueryRowContext(ctx, query, f.UsuarioID, f.MaterialID, f.Pergunta, f.Resposta).
		Scan(&f.ID, &f.ProximaRevisao, &f.DataCriacao)
}

func (r *EstudoPostgres) ListarFlashcards(ctx context.Context, usuarioID int, materialID int) ([]estudo.Flashcard, error) {
	query := `SELECT id, usuario_id, material_id, pergunta, resposta, dificuldade, proxima_revisao, data_criacao
	          FROM flashcards WHERE usuario_id = $1`
	args := []interface{}{usuarioID}
	if materialID > 0 {
		query += " AND material_id = $2"
		args = append(args, materialID)
	}
	query += " ORDER BY proxima_revisao ASC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var flashcards []estudo.Flashcard
	for rows.Next() {
		var f estudo.Flashcard
		if err := rows.Scan(&f.ID, &f.UsuarioID, &f.MaterialID, &f.Pergunta, &f.Resposta, &f.Dificuldade, &f.ProximaRevisao, &f.DataCriacao); err != nil {
			return nil, err
		}
		flashcards = append(flashcards, f)
	}
	return flashcards, nil
}

func (r *EstudoPostgres) AtualizarDificuldade(ctx context.Context, id int, dificuldade int) error {
	// Simple SRS logic: harder items repeat sooner
	daysToAdd := 1
	switch dificuldade {
	case 1: // Fácil
		daysToAdd = 7
	case 2: // Médio
		daysToAdd = 3
	default: // Difícil
		daysToAdd = 1
	}

	nextReview := time.Now().AddDate(0, 0, daysToAdd)
	query := `UPDATE flashcards SET dificuldade = $1, proxima_revisao = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, dificuldade, nextReview, id)
	return err
}

func (r *EstudoPostgres) DeletarFlashcard(ctx context.Context, id, usuarioID int) error {
	query := `DELETE FROM flashcards WHERE id = $1 AND usuario_id = $2`
	_, err := r.DB.ExecContext(ctx, query, id, usuarioID)
	return err
}
