package material

type Repository interface {
	Listar(limit, offset int) ([]Material, error)
	BuscarPorID(id int) (*Material, error)
	Pesquisar(termo string, categoria string, tags []string, limit, offset int) ([]Material, error)
	Criar(m *Material) error
	Atualizar(m *Material) error
	Deletar(id int) error

	// Empréstimos
	SalvarEmprestimo(e *Emprestimo) error
	ListarEmprestimosPorUsuario(usuarioID int) ([]Emprestimo, error)

	// Interações
	SalvarAvaliacao(a *Avaliacao) error
	ListarAvaliacoesPorMaterial(materialID int) ([]Avaliacao, error)
	AdicionarFavorito(f *Favorito) error
	RemoverFavorito(usuarioID, materialID int) error
	ListarFavoritosPorUsuario(usuarioID int) ([]Material, error)
	RegistrarLeitura(h *HistoricoLeitura) error
	ListarHistoricoPorUsuario(usuarioID int) ([]Material, error)
	ObterRecomendacoes(usuarioID int, limit int) ([]Material, error)
}
