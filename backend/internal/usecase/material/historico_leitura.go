package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"time"
)

type HistoricoLeituraUseCase struct {
	Repo material.Repository
}

func (uc *HistoricoLeituraUseCase) Execute(usuarioID, materialID int) error {
	h := &material.HistoricoLeitura{
		UsuarioID:  usuarioID,
		MaterialID: materialID,
		Data:       time.Now(),
	}
	return uc.Repo.RegistrarLeitura(h)
}

func (uc *HistoricoLeituraUseCase) Listar(usuarioID int) ([]material.Material, error) {
	return uc.Repo.ListarHistoricoPorUsuario(usuarioID)
}
