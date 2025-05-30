package usuario

import (
	"context"
	"biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/pkg/hash"
)

type CadastrarUsuarioUseCase struct {
	repo usuario.UsuarioRepository
}

func NewCadastrarUsuario(repo usuario.UsuarioRepository) *CadastrarUsuarioUseCase {
	return &CadastrarUsuarioUseCase{repo}
}

func (uc *CadastrarUsuarioUseCase) Execute(ctx context.Context, u usuario.Usuario) error {
	u.Senha = hash.GerarHash(u.Senha)
	return uc.repo.Salvar(ctx, u)
}
