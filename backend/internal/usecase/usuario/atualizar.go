package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"context"
)

type AtualizarUsuarioUseCase struct {
	repo usuario.UsuarioRepository
}

func NewAtualizarUsuario(repo usuario.UsuarioRepository) *AtualizarUsuarioUseCase {
	return &AtualizarUsuarioUseCase{repo}
}

func (uc *AtualizarUsuarioUseCase) Execute(ctx context.Context, u usuario.Usuario) error {
	return uc.repo.Atualizar(ctx, &u)
}
