package estudo

import (
	"context"
)

type Repository interface {
	// Anotações
	CriarAnotacao(ctx context.Context, a *Anotacao) error
	ListarAnotacoes(ctx context.Context, usuarioID, materialID int) ([]Anotacao, error)
	DeletarAnotacao(ctx context.Context, id, usuarioID int) error

	// Flashcards
	CriarFlashcard(ctx context.Context, f *Flashcard) error
	ListarFlashcards(ctx context.Context, usuarioID int, materialID int) ([]Flashcard, error)
	AtualizarDificuldade(ctx context.Context, id int, dificuldade int) error
	DeletarFlashcard(ctx context.Context, id, usuarioID int) error
}
