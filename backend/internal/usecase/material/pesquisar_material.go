package material

import "biblioteca-digital-api/internal/domain/material"

type PesquisarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *PesquisarMaterialUseCase) Execute(termo, categoria string, tags []string, limit, offset int) ([]material.Material, error) {
	return uc.Repo.Pesquisar(termo, categoria, tags, limit, offset)
}
