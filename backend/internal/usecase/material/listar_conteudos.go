package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"context"
	"fmt"
	"time"
)

type ListarConteudosUseCase struct {
	Repo      material.Repository
	Harvester Harvester
	Cache     cache.Cache
}

func (uc *ListarConteudosUseCase) Execute(ctx context.Context, limit, offset int) ([]material.Material, error) {
	cacheKey := fmt.Sprintf("list:%d:%d", limit, offset)
	if uc.Cache != nil {
		var cached []material.Material
		if found := uc.Cache.Get(cacheKey, &cached); found {
			return cached, nil
		}
	}

	materiais, err := uc.Repo.Listar(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Se estiver vazio e for a primeira página, tenta buscar algo do harvester para não ficar vazio
	if len(materiais) == 0 && offset == 0 && uc.Harvester != nil {
		externalMats, err := uc.Harvester.Search(ctx, "", "", "", 0, 0, limit)
		if err == nil {
			for i := range externalMats {
				m := externalMats[i]
				_ = uc.Repo.Criar(ctx, &m)
				materiais = append(materiais, m)
			}
		}
	}

	if uc.Cache != nil {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
