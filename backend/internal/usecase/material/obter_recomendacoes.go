package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
)

type ObterRecomendacoesUseCase struct {
	Repo material.Repository
}

func (uc *ObterRecomendacoesUseCase) Execute(ctx context.Context, usuarioID, limit int) ([]material.Material, error) {
	return uc.Repo.ObterRecomendacoes(ctx, usuarioID, limit)
}
