package usuario

import "context"

type UsuarioRepository interface {
	Salvar(ctx context.Context, usuario *Usuario) error
	BuscarPorEmail(ctx context.Context, email string) (*Usuario, error)
	ListarInteresses(ctx context.Context, id int) ([]string, error)
	AtualizarSenha(ctx context.Context, email string, novaSenha string) error
	Atualizar(ctx context.Context, u *Usuario) error
}
