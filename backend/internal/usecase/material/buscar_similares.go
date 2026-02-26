package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
)

type BuscarSimilaresUseCase struct {
	Repo material.Repository
}

func (uc *BuscarSimilaresUseCase) Execute(ctx context.Context, materialID int, limit int) ([]material.Material, error) {
	return uc.Repo.BuscarSimilares(ctx, materialID, limit)
}
