package material

import "context"

type Repository interface {
	Listar(ctx context.Context, limit, offset int) ([]Material, error)
	BuscarPorID(ctx context.Context, id int) (*Material, error)
	Pesquisar(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]Material, error)
	Sugerir(ctx context.Context, termo string, limit int) ([]Material, error)
	BuscarSimilares(ctx context.Context, materialID int, limit int) ([]Material, error)
	Criar(ctx context.Context, m *Material) error
	Atualizar(ctx context.Context, m *Material) error
	Deletar(ctx context.Context, id int) error

	// Interações
	AdicionarFavorito(ctx context.Context, f *Favorito) error
	RemoverFavorito(ctx context.Context, usuarioID, materialID int) error
	ListarFavoritosPorUsuario(ctx context.Context, usuarioID int) ([]Material, error)
	RegistrarLeitura(ctx context.Context, h *HistoricoLeitura) error
	ListarHistoricoPorUsuario(ctx context.Context, usuarioID int) ([]Material, error)
}
