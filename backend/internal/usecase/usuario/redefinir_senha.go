package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/pkg/hash"
	"context"
	"errors"
)

type RedefinirSenhaUseCase struct {
	repo usuario.UsuarioRepository
}

func NewRedefinirSenhaUseCase(repo usuario.UsuarioRepository) *RedefinirSenhaUseCase {
	return &RedefinirSenhaUseCase{repo}
}

func (uc *RedefinirSenhaUseCase) Execute(ctx context.Context, email, novaSenha string) error {
	// 1. Verificar se o usuário existe
	user, err := uc.repo.BuscarPorEmail(ctx, email)
	if err != nil || user == nil {
		return errors.New("usuário não encontrado")
	}

	// 2. Gerar hash da nova senha
	hashedPassword, err := hash.GerarHash(novaSenha)
	if err != nil {
		return err
	}

	// 3. Atualizar no repositório
	return uc.repo.AtualizarSenha(ctx, email, hashedPassword)
}
