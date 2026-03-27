package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/utils"
	"context"
	"fmt"
	"strings"
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

	// 1. Local Search (Database)
	// We oversample to compensate for any later filtering
	fetchLimit := limit + (limit / 2)
	materiaisIniciais, localErr := uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, fetchLimit, offset, sort)
	if localErr != nil {
		return nil, localErr
	}

	// 2. High Density Harvesting (Trigger if local database is too sparse)
	// We define "sparse" as having less than half of the requested limit OR if the user specifically asked for a high limit
	isSparse := len(materiaisIniciais) < limit || limit > 50 
	hasQuery := termo != "" || categoria != ""

	if hasQuery && isSparse && uc.Harvester != nil {
		// Calculate harvest depth based on the machine's strength (implied by requester)
		harvestLimit := limit
		if limit < 40 { harvestLimit = 60 } // Minimum decent harvest

		harvested, err := uc.Harvester.Search(ctx, termo, categoria, fonte, anoInicio, anoFim, harvestLimit)
		if err == nil && len(harvested) > 0 {
			// Persist in background (non-blocking for results but we want them for next time)
			for i := range harvested {
				go func(m *material.Material) {
					// We use a separate context for persistence to avoid truncation if response is sent
					_ = uc.Repo.Criar(context.Background(), m)
				}(&harvested[i])
			}
			// Combine results
			materiaisIniciais = append(materiaisIniciais, harvested...)
		}
	}

	// 3. Deduplication and Verification
	var materiais []material.Material
	uniqueMap := make(map[string]material.Material)
	urlsToCheck := make([]string, 0)

	for _, m := range materiaisIniciais {
		sig := m.ExternoID
		if sig == "" {
			sig = strings.ToLower(m.Titulo + ":" + m.Autor)
		}
		if _, exists := uniqueMap[sig]; !exists {
			uniqueMap[sig] = m
			if m.PDFURL != "" {
				urlsToCheck = append(urlsToCheck, m.PDFURL)
			}
		}
	}

	// Parallel Batch Verification
	statusMap := make(map[string]bool)
	if uc.Verifier != nil && len(urlsToCheck) > 0 {
		verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second) // Generous timeout for high volume
		defer cancel()
		statusMap = uc.Verifier.VerifyBatch(verifyCtx, urlsToCheck)
	}

	for _, m := range uniqueMap {
		if m.PDFURL != "" {
			// Optimistic: Only exclude if explicitly reported as DEAD (404/403 etc)
			if alive, exists := statusMap[m.PDFURL]; exists && !alive {
				continue
			}
		}
		materiais = append(materiais, m)
	}

	// 4. Force limit and final sort (optional but good for UI consistency)
	if limit > 0 && len(materiais) > limit {
		materiais = materiais[:limit]
	}

	// 5. Caching
	if uc.Cache != nil && len(materiais) > 0 {
		uc.Cache.Set(cacheKey, materiais, 10*time.Minute)
	}

	return materiais, nil
}
