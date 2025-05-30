package usuario

type Repository interface {
	Salvar(usuario *Usuario) error
	BuscarPorEmail(email string) (*Usuario, error)
	ListarInteresses(id int) ([]string, error)
}

