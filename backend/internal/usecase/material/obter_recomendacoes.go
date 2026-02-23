package material

import "biblioteca-digital-api/internal/domain/material"

type ObterRecomendacoesUseCase struct {
	Repo material.Repository
}

func (uc *ObterRecomendacoesUseCase) Execute(usuarioID int, limit int) ([]material.Material, error) {
	return uc.Repo.ObterRecomendacoes(usuarioID, limit)
}
