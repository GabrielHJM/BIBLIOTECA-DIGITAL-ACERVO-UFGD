package usuario

import "biblioteca-digital-api/internal/domain/usuario"

type ListarInteressesUseCase struct {
	Repo usuario.Repository
}

func (uc *ListarInteressesUseCase) Execute(id int) ([]string, error) {
	return uc.Repo.ListarInteresses(id)
}
