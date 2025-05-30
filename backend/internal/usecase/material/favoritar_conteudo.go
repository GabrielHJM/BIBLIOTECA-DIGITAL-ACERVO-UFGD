package material

import "biblioteca-digital-api/internal/domain/material"

type FavoritarConteudoUseCase struct {
	Repo material.Repository
}

func (uc *FavoritarConteudoUseCase) Execute(id int) error {
	return uc.Repo.Favoritar(id)
}
