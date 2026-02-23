package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/pkg/hash"
	"context"
)

type CadastrarUsuarioUseCase struct {
	repo usuario.UsuarioRepository
}

func NewCadastrarUsuario(repo usuario.UsuarioRepository) *CadastrarUsuarioUseCase {
	return &CadastrarUsuarioUseCase{repo}
}

func (uc *CadastrarUsuarioUseCase) Execute(ctx context.Context, u usuario.Usuario) error {
	hashedPassword, err := hash.GerarHash(u.Senha)
	if err != nil {
		return err
	}
	u.Senha = hashedPassword
	return uc.repo.Salvar(ctx, &u)
}
