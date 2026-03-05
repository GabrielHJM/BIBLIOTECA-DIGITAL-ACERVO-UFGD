package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
)

type BuscarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *BuscarMaterialUseCase) Execute(ctx context.Context, id int) (*material.Material, error) {
	m, err := uc.Repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	return m, nil
}
