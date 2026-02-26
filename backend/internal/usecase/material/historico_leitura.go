package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"time"
)

type HistoricoLeituraUseCase struct {
	Repo material.Repository
}

func (uc *HistoricoLeituraUseCase) Execute(ctx context.Context, usuarioID, materialID int) error {
	h := &material.HistoricoLeitura{
		UsuarioID:  usuarioID,
		MaterialID: materialID,
		Data:       time.Now(),
	}
	return uc.Repo.RegistrarLeitura(ctx, h)
}

func (uc *HistoricoLeituraUseCase) Listar(ctx context.Context, usuarioID int) ([]material.Material, error) {
	return uc.Repo.ListarHistoricoPorUsuario(ctx, usuarioID)
}
