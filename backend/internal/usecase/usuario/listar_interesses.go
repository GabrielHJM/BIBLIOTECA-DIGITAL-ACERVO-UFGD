package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"context"
)

type ListarInteressesUseCase struct {
	Repo usuario.UsuarioRepository
}

func (uc *ListarInteressesUseCase) Execute(ctx context.Context, id int) ([]string, error) {
	return uc.Repo.ListarInteresses(ctx, id)
}
