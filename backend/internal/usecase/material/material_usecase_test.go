package material

import (
	"context"
	"errors"
	"testing"
	"time"

	"biblioteca-digital-api/internal/domain/material"
)

// MockRepository implements material.Repository
type MockRepository struct {
	material.Repository
	ListarFunc    func(ctx context.Context, limit, offset int) ([]material.Material, error)
	PesquisarFunc func(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error)
	CriarFunc     func(ctx context.Context, m *material.Material) error
}

func (m *MockRepository) Listar(ctx context.Context, limit, offset int) ([]material.Material, error) {
	return m.ListarFunc(ctx, limit, offset)
}

func (m *MockRepository) Pesquisar(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
	return m.PesquisarFunc(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, limit, offset, sort)
}

func (m *MockRepository) Criar(ctx context.Context, mat *material.Material) error {
	if m.CriarFunc != nil {
		return m.CriarFunc(ctx, mat)
	}
	return nil
}

// MockHarvester implements Harvester interface
type MockHarvester struct {
	SearchFunc func(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error)
}

func (m *MockHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	return m.SearchFunc(ctx, query, category, source, startYear, endYear, limit)
}

// MockCache implements Cache interface
type MockCache struct {
	GetFunc func(key string) (interface{}, bool)
	SetFunc func(key string, value interface{}, duration time.Duration)
}

func (m *MockCache) Get(key string) (interface{}, bool) {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return nil, false
}

func (m *MockCache) Set(key string, value interface{}, duration time.Duration) {
	if m.SetFunc != nil {
		m.SetFunc(key, value, duration)
	}
}

func TestListarConteudosExecute(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso_busca_local", func(t *testing.T) {
		expectedMats := []material.Material{{ID: 1, Titulo: "Livro 1"}}
		repo := &MockRepository{
			ListarFunc: func(ctx context.Context, limit, offset int) ([]material.Material, error) {
				return expectedMats, nil
			},
		}
		uc := &ListarConteudosUseCase{Repo: repo}

		res, err := uc.Execute(ctx, 10, 0)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(res) != 1 || res[0].Titulo != "Livro 1" {
			t.Errorf("resultado inesperado: %v", res)
		}
	})

	t.Run("sucesso_busca_harvester_vazio_local", func(t *testing.T) {
		repo := &MockRepository{
			ListarFunc: func(ctx context.Context, limit, offset int) ([]material.Material, error) {
				return []material.Material{}, nil
			},
		}
		harvester := &MockHarvester{
			SearchFunc: func(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
				return []material.Material{{Titulo: "Externo 1"}}, nil
			},
		}
		uc := &ListarConteudosUseCase{Repo: repo, Harvester: harvester}

		res, err := uc.Execute(ctx, 10, 0)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(res) != 1 || res[0].Titulo != "Externo 1" {
			t.Errorf("resultado inesperado: %v", res)
		}
	})

	t.Run("erro_repositorio", func(t *testing.T) {
		repo := &MockRepository{
			ListarFunc: func(ctx context.Context, limit, offset int) ([]material.Material, error) {
				return nil, errors.New("erro de bd")
			},
		}
		uc := &ListarConteudosUseCase{Repo: repo}

		_, err := uc.Execute(ctx, 10, 0)
		if err == nil {
			t.Fatal("esperava erro, recebeu nil")
		}
	})
}

func TestPesquisarMaterialExecute(t *testing.T) {
	ctx := context.Background()

	t.Run("sucesso_busca_local_apenas", func(t *testing.T) {
		repo := &MockRepository{
			PesquisarFunc: func(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
				return []material.Material{{ID: 1, Titulo: "Local 1"}}, nil
			},
		}
		uc := &PesquisarMaterialUseCase{Repo: repo}

		res, err := uc.Execute(ctx, "termo", "", "", 0, 0, nil, 10, 0, "")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if len(res) != 1 || res[0].Titulo != "Local 1" {
			t.Errorf("resultado inesperado: %v", res)
		}
	})

	t.Run("sucesso_mesclagem_externo", func(t *testing.T) {
		repo := &MockRepository{
			PesquisarFunc: func(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
				return []material.Material{{ID: 1, Titulo: "Local 1", ExternoID: "ext_1"}}, nil
			},
			CriarFunc: func(ctx context.Context, m *material.Material) error {
				m.ID = 2
				return nil
			},
		}
		harvester := &MockHarvester{
			SearchFunc: func(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
				return []material.Material{
					{Titulo: "Local 1", ExternoID: "ext_1"}, // Duplicado
					{Titulo: "Novo Externo", ExternoID: "ext_2"},
				}, nil
			},
		}
		uc := &PesquisarMaterialUseCase{Repo: repo, Harvester: harvester}

		res, err := uc.Execute(ctx, "termo", "", "", 0, 0, nil, 10, 0, "")
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		// Deve ter 2 resultados: Local 1 (do repo) e Novo Externo (do harvester)
		if len(res) != 2 {
			t.Errorf("esperava 2 resultados, recebeu %d", len(res))
		}
	})
}

// End of tests
