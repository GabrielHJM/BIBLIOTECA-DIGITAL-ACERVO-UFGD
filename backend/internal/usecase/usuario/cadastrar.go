package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"biblioteca-digital-api/pkg/hash"
	"context"
	"errors"
	"regexp"
)

type CadastrarUsuarioUseCase struct {
	repo usuario.UsuarioRepository
}

func NewCadastrarUsuario(repo usuario.UsuarioRepository) *CadastrarUsuarioUseCase {
	return &CadastrarUsuarioUseCase{repo}
}

func (uc *CadastrarUsuarioUseCase) Execute(ctx context.Context, u usuario.Usuario) error {
	// Validação simples de email
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(u.Email) {
		return errors.New("formato de email inválido")
	}

	hashedPassword, err := hash.GerarHash(u.Senha)
	if err != nil {
		return err
	}
	u.Senha = hashedPassword
	return uc.repo.Salvar(ctx, &u)
}
