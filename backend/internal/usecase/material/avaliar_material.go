package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"time"
)

type AvaliarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *AvaliarMaterialUseCase) Execute(usuarioID, materialID, nota int, comentario string) error {
	a := &material.Avaliacao{
		UsuarioID:  usuarioID,
		MaterialID: materialID,
		Nota:       nota,
		Comentario: comentario,
		Data:       time.Now(),
	}
	return uc.Repo.SalvarAvaliacao(a)
}
