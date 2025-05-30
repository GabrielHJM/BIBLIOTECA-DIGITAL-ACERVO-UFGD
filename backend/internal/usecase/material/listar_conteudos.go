package material

import "biblioteca-digital-api/internal/domain/material"

type ListarConteudosUseCase struct {
	Repo material.Repository
}

func (uc *ListarConteudosUseCase) Execute() ([]material.Material, error) {
	return uc.Repo.Listar()
}
