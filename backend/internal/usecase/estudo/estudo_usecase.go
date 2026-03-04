package estudo

import (
	"biblioteca-digital-api/internal/domain/estudo"
	"biblioteca-digital-api/internal/domain/material"
	"context"
)

type UseCase struct {
	repo         estudo.Repository
	materialRepo material.Repository
}

func NewUseCase(repo estudo.Repository, materialRepo material.Repository) *UseCase {
	return &UseCase{repo, materialRepo}
}

// Flashcards
func (uc *UseCase) CriarFlashcard(ctx context.Context, f *estudo.Flashcard) error {
	return uc.repo.CriarFlashcard(ctx, f)
}

func (uc *UseCase) ListarFlashcards(ctx context.Context, usuarioID int, materialID int) ([]estudo.Flashcard, error) {
	return uc.repo.ListarFlashcards(ctx, usuarioID, materialID)
}

func (uc *UseCase) AtualizarDificuldade(ctx context.Context, id int, dificuldade int) error {
	return uc.repo.AtualizarDificuldade(ctx, id, dificuldade)
}
