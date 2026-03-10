package material

import (
	"biblioteca-digital-api/internal/repository"
	"context"
)

type AvaliarMaterialUseCase struct {
	Repo *repository.MaterialPostgres
}

func (uc *AvaliarMaterialUseCase) Execute(ctx context.Context, usuarioID, materialID int, nota float64) error {
	return uc.Repo.Avaliar(ctx, usuarioID, materialID, nota)
}
