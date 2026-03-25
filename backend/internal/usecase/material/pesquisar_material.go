package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/utils"
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
	Verifier  *utils.URLVerifier
}

func (uc *PesquisarMaterialUseCase) Execute(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
	cacheKey := fmt.Sprintf("search:%s:%s:%s:%d:%d:%d:%d:%s", termo, categoria, fonte, anoInicio, anoFim, limit, offset, sort)
	if uc.Cache != nil && sort != "random" {
		var cached []material.Material
		if found := uc.Cache.Get(cacheKey, &cached); found {
			return cached, nil
		}
	}

	// 1. Oversample from DB to compensate for potential filtered items
	fetchLimit := limit
	if limit > 0 {
		fetchLimit = limit * 2
	}

	// Busca local (Banco de Dados)
	materiaisIniciais, localErr := uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, fetchLimit, offset, sort)
	if localErr != nil {
		return nil, fmt.Errorf("erro na busca local: %w", localErr)
	}

	// 2. Harvester Discovery (Discover new materials from external APIs)
	// We trigger this if searching for a term or specific category
	if (termo != "" || categoria != "") && uc.Harvester != nil {
		harvestLimit := 40 // Increased yield per source
		harvested, err := uc.Harvester.Search(ctx, termo, categoria, fonte, anoInicio, anoFim, harvestLimit)
		if err == nil && len(harvested) > 0 {
			// Persist harvested results to DB (in background or simple sequence)
			for i := range harvested {
				_ = uc.Repo.Criar(ctx, &harvested[i])
			}
			
			// If local results are sparse, merge harvested ones
			if len(materiaisIniciais) < limit {
				materiaisIniciais = append(materiaisIniciais, harvested...)
			}
		}
	}

	var materiais []material.Material

	// 3. Intelligent Verification Algorithm
	if len(materiaisIniciais) > 0 && uc.Verifier != nil {
		var urlsToCheck []string
		
		// Deduplicate and filter by relevance/quality if multiple sources
		uniqueMap := make(map[string]material.Material)
		for _, m := range materiaisIniciais {
			uID := m.ExternoID
			if uID == "" {
				uID = m.Titulo + m.Autor
			}
			if _, exists := uniqueMap[uID]; !exists {
				uniqueMap[uID] = m
				if m.PDFURL != "" {
					urlsToCheck = append(urlsToCheck, m.PDFURL)
				}
			}
		}

		// Perform parallel check with strict timeout context
		verifyCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond) // Slightly more time for batch check
		defer cancel()
		
		statusMap := uc.Verifier.VerifyBatch(verifyCtx, urlsToCheck)

		var filtered []material.Material
		for _, m := range uniqueMap {
			if m.PDFURL != "" {
				if alive, exists := statusMap[m.PDFURL]; exists && !alive {
					continue // Skip dead links
				}
			}
			filtered = append(filtered, m)
		}
		
		// Result limit enforcement after filtering
		if limit > 0 && len(filtered) > limit {
			filtered = filtered[:limit]
		}
		materiais = filtered
	} else {
		materiais = materiaisIniciais
	}

	if uc.Cache != nil && sort != "random" && len(materiais) > 0 {
		uc.Cache.Set(cacheKey, materiais, 10*time.Minute)
	}

	return materiais, nil
}
