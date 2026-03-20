package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"context"
	"fmt"
	"time"
)

type Harvester interface {
	Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error)
}

type PesquisarMaterialUseCase struct {
	Repo      material.Repository
	Harvester Harvester
	Cache     cache.Cache
}

func (uc *PesquisarMaterialUseCase) Execute(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
	cacheKey := fmt.Sprintf("search:%s:%s:%s:%d:%d:%d:%d:%s", termo, categoria, fonte, anoInicio, anoFim, limit, offset, sort)
	if uc.Cache != nil && sort != "random" {
		var cached []material.Material
		if found := uc.Cache.Get(cacheKey, &cached); found {
			return cached, nil
		}
	}

	var (
		materiais []material.Material
		localErr  error
	)

	// Busca local (Banco de Dados) executada de forma nativa e ultra rápida
	materiais, localErr = uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, limit, offset, sort)

	if localErr != nil {
		return nil, fmt.Errorf("erro na busca local: %w", localErr)
	}

	if uc.Cache != nil && sort != "random" && len(materiais) > 0 {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
