package usuario

type Service interface {
	Cadastrar(usuario *Usuario) error
	Login(email, senha string) (*Usuario, error)
	ListarInteresses(id int) ([]string, error)
}
