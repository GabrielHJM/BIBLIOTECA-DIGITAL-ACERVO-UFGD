package usuario

import (
	"context"
	"errors"
	"biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/pkg/auth"
	"biblioteca-digital-api/pkg/hash"
)

type LoginUseCase struct {
	repo usuario.UsuarioRepository
}

func NewLoginUseCase(repo usuario.UsuarioRepository) *LoginUseCase {
	return &LoginUseCase{repo}
}

func (uc *LoginUseCase) Execute(ctx context.Context, email, senha string) (string, error) {
	user, err := uc.repo.BuscarPorEmail(ctx, email)
	if err != nil || !hash.VerificarHash(senha, user.Senha) {
		return "", errors.New("credenciais inválidas")
	}
	return auth.GerarToken(user.ID)
}
