package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"context"
	"fmt"
	"sync"
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
		materiais    []material.Material
		externalMats []material.Material
		localErr     error
		extErr       error
		wg           sync.WaitGroup
	)

	// Busca local (Banco de Dados)
	wg.Add(1)
	go func() {
		defer wg.Done()
		materiais, localErr = uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, limit, offset, sort)
	}()

	// Só busca externo se:
	// 1. Harvester estiver configurado
	// 2. For a primeira página (offset == 0)
	// 3. Não for ordem aleatória (aleatório é local)
	// Aguardamos a busca local terminar primeiro para decidir se precisamos da externa (Otimização de Latência)
	wg.Wait()

	if localErr != nil {
		return nil, fmt.Errorf("erro na busca local: %w", localErr)
	}

	// Se temos poucos resultados locais e estamos na primeira página, disparamos a busca externa
	if uc.Harvester != nil && offset == 0 && len(materiais) < limit && sort != "random" {
		// Timeout estrito de 3 segundos para buscas externas (Fast Fail)
		harvestCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		externalMats, extErr = uc.Harvester.Search(harvestCtx, termo, categoria, fonte, anoInicio, anoFim, limit)

		// Mesclar resultados externos se houver
		if extErr == nil && len(externalMats) > 0 {
			seen := make(map[string]bool)
			for _, m := range materiais {
				if m.ExternoID != "" {
					seen[m.ExternoID] = true
				}
			}

			newMats := []material.Material{}
			for i := range externalMats {
				em := &externalMats[i]
				if em.ExternoID != "" && seen[em.ExternoID] {
					continue
				}

				// Salvamento ASSÍNCRONO no banco (Background) para não custar latência ao usuário
				go func(mat material.Material) {
					bgCtx, bgCancel := context.WithTimeout(context.Background(), 10*time.Second)
					defer bgCancel()
					_ = uc.Repo.Criar(bgCtx, &mat)
				}(*em)

				if em.ExternoID != "" {
					seen[em.ExternoID] = true
				}

				if len(materiais)+len(newMats) < limit {
					newMats = append(newMats, *em)
				}
			}
			materiais = append(materiais, newMats...)
		}
	}

	if uc.Cache != nil && sort != "random" && len(materiais) > 0 {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
