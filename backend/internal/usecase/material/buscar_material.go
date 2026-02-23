package material

import "biblioteca-digital-api/internal/domain/material"

type BuscarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *BuscarMaterialUseCase) Execute(id int) (*material.Material, error) {
	return uc.Repo.BuscarPorID(id)
}
