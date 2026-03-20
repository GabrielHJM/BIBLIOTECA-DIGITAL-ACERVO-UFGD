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

	// As bucas externas na API síncrona foram removidas de propósito
	// para evitar extrema latência e 429 Too Many Requests.
	// O background worker populará esse acervo automaticamente.

	if uc.Cache != nil {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
