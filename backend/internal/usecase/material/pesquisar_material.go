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

	wg.Add(1)
	go func() {
		defer wg.Done()
		materiais, localErr = uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, limit, offset, sort)
	}()

	// Só busca externo se não for uma busca com offset (paginação profunda geralmente é baseada em banco local)
	// E se tivermos um Harvester configurado
	if uc.Harvester != nil && offset == 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Busca externa pode ser lenta, usamos um timeout específico se necessário ou o próprio ctx
			externalMats, extErr = uc.Harvester.Search(ctx, termo, categoria, fonte, anoInicio, anoFim, limit)
		}()
	}

	wg.Wait()

	if localErr != nil {
		return nil, fmt.Errorf("erro na busca local: %w", localErr)
	}

	// Mesclar resultados
	if extErr == nil && len(externalMats) > 0 {
		// Mapa para evitar duplicados por ExternoID
		seen := make(map[string]bool)
		for _, m := range materiais {
			if m.ExternoID != "" {
				seen[m.ExternoID] = true
			}
		}

		for i := range externalMats {
			em := &externalMats[i]
			if em.ExternoID != "" && seen[em.ExternoID] {
				continue
			}

			// Salva resultados externos no banco para gerar o ID sincornamente antes de retornar ao usuário
			err := uc.Repo.Criar(ctx, em)
			if err != nil {
				fmt.Printf("Erro ao salvar material externo: %v\n", err)
				continue
			}

			if em.ExternoID != "" {
				seen[em.ExternoID] = true
			}

			if len(materiais) < limit {
				materiais = append(materiais, *em)
			}
		}
	}

	if uc.Cache != nil && sort != "random" {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
