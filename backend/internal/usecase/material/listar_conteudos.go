package material

import "biblioteca-digital-api/internal/domain/material"

type ListarConteudosUseCase struct {
	Repo material.Repository
}

func (uc *ListarConteudosUseCase) Execute(limit, offset int) ([]material.Material, error) {
	return uc.Repo.Listar(limit, offset)
}
