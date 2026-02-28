package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/ai"
	"context"
	"fmt"
)

type AIUseCase struct {
	repo   material.Repository
	gemini *ai.GeminiClient
}

func NewAIUseCase(repo material.Repository, gemini *ai.GeminiClient) *AIUseCase {
	return &AIUseCase{repo, gemini}
}

func (uc *AIUseCase) PerguntarLivro(ctx context.Context, materialID int, pergunta string) (string, error) {
	m, err := uc.repo.BuscarPorID(ctx, materialID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Você é um assistente virtual de uma biblioteca digital.
Responda a perguntas sobre o seguinte livro:
Título: %s
Autor: %s
Categoria: %s
Descrição/Resumo: %s

Pergunta do Usuário: %s

Responda de forma clara, educada e baseada estritamente nas informações fornecidas ou no conhecimento geral sobre esta obra.`,
		m.Titulo, m.Autor, m.Categoria, m.Descricao, pergunta)

	return uc.gemini.GenerateContent(prompt)
}

func (uc *AIUseCase) GerarResumo(ctx context.Context, materialID int) (string, error) {
	m, err := uc.repo.BuscarPorID(ctx, materialID)
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`Crie um resumo executivo e acadêmico para o livro abaixo, destacando os pontos principais, objetivos e público-alvo.
Use tópicos (bullet points) se necessário.

Título: %s
Autor: %s
Categoria: %s
Descrição Original: %s

Gere um texto premium para estudantes universitários.`, m.Titulo, m.Autor, m.Categoria, m.Descricao)

	return uc.gemini.GenerateContent(prompt)
}
