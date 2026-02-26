package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"time"
)

type AvaliarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *AvaliarMaterialUseCase) Execute(ctx context.Context, usuarioID, materialID, nota int, comentario string) error {
	a := &material.Avaliacao{
		UsuarioID:  usuarioID,
		MaterialID: materialID,
		Nota:       nota,
		Comentario: comentario,
		Data:       time.Now(),
	}

	return uc.Repo.SalvarAvaliacao(ctx, a)
}

func (uc *AvaliarMaterialUseCase) Listar(ctx context.Context, materialID int) ([]material.Avaliacao, error) {
	return uc.Repo.ListarAvaliacoesPorMaterial(ctx, materialID)
}
