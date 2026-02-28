package estudo

import (
	"biblioteca-digital-api/internal/domain/estudo"
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/ai"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type UseCase struct {
	repo         estudo.Repository
	materialRepo material.Repository
	gemini       *ai.GeminiClient
}

func NewUseCase(repo estudo.Repository, materialRepo material.Repository, gemini *ai.GeminiClient) *UseCase {
	return &UseCase{repo, materialRepo, gemini}
}

// Anotações
func (uc *UseCase) CriarAnotacao(ctx context.Context, a *estudo.Anotacao) error {
	return uc.repo.CriarAnotacao(ctx, a)
}

func (uc *UseCase) ListarAnotacoes(ctx context.Context, usuarioID, materialID int) ([]estudo.Anotacao, error) {
	return uc.repo.ListarAnotacoes(ctx, usuarioID, materialID)
}

func (uc *UseCase) DeletarAnotacao(ctx context.Context, id, usuarioID int) error {
	return uc.repo.DeletarAnotacao(ctx, id, usuarioID)
}

// Flashcards
func (uc *UseCase) CriarFlashcard(ctx context.Context, f *estudo.Flashcard) error {
	return uc.repo.CriarFlashcard(ctx, f)
}

func (uc *UseCase) ListarFlashcards(ctx context.Context, usuarioID int, materialID int) ([]estudo.Flashcard, error) {
	return uc.repo.ListarFlashcards(ctx, usuarioID, materialID)
}

func (uc *UseCase) GerarFlashcardsIA(ctx context.Context, usuarioID, materialID int) ([]estudo.Flashcard, error) {
	m, err := uc.materialRepo.BuscarPorID(ctx, materialID)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`Com base no livro abaixo, gere 5 flashcards (perguntas e respostas curtas) para estudo.
Responda APENAS com um array JSON no formato: [{"pergunta": "...", "resposta": "..."}]

Título: %s
Autor: %s
Descrição: %s`, m.Titulo, m.Autor, m.Descricao)

	resp, err := uc.gemini.GenerateContent(prompt)
	if err != nil {
		return nil, err
	}

	// Limpeza básica do JSON em caso de markdown
	resp = strings.TrimPrefix(resp, "```json")
	resp = strings.TrimSuffix(resp, "```")
	resp = strings.TrimSpace(resp)

	var items []struct {
		Pergunta string `json:"pergunta"`
		Resposta string `json:"resposta"`
	}
	if err := json.Unmarshal([]byte(resp), &items); err != nil {
		return nil, fmt.Errorf("erro ao processar flashcards da IA: %v", err)
	}

	var flashcards []estudo.Flashcard
	for _, item := range items {
		f := estudo.Flashcard{
			UsuarioID:  usuarioID,
			MaterialID: materialID,
			Pergunta:   item.Pergunta,
			Resposta:   item.Resposta,
		}
		if err := uc.repo.CriarFlashcard(ctx, &f); err == nil {
			flashcards = append(flashcards, f)
		}
	}

	return flashcards, nil
}

func (uc *UseCase) AtualizarDificuldade(ctx context.Context, id int, dificuldade int) error {
	return uc.repo.AtualizarDificuldade(ctx, id, dificuldade)
}
