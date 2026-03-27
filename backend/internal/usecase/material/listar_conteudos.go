package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/utils"
	"context"
	"fmt"
	"time"
)

type ListarConteudosUseCase struct {
	Repo      material.Repository
	Harvester Harvester
	Cache     cache.Cache
	Verifier  *utils.URLVerifier
}

func (uc *ListarConteudosUseCase) Execute(ctx context.Context, limit, offset int) ([]material.Material, error) {
	cacheKey := fmt.Sprintf("list:%d:%d", limit, offset)
	if uc.Cache != nil {
		var cached []material.Material
		if found := uc.Cache.Get(cacheKey, &cached); found {
			return cached, nil
		}
	}

	materiaisIniciais, err := uc.Repo.Listar(ctx, limit*2, offset)
	if err != nil {
		return nil, err
	}

	// Real-time Infinite Harvester Algorithm:
	// If local results are low, we trigger a light-weight synchronous harvest
	// BEFORE the final result build to ensure immediate visibility.
	if len(materiaisIniciais) < limit && uc.Harvester != nil {
		harvestLimit := limit * 2
		discoveryCats := []string{"TECNOLOGIA", "CIÊNCIAS", "EDUCAÇÃO", "SAÚDE", "DIREITO"}
		cat := discoveryCats[time.Now().UnixNano()%int64(len(discoveryCats))]
		
		harvested, err := uc.Harvester.Search(ctx, "", cat, "", 0, 0, harvestLimit)
		if err == nil && len(harvested) > 0 {
			for i := range harvested {
				_ = uc.Repo.Criar(ctx, &harvested[i])
			}
			// Important: Merge newly found materials so they show up in THIS request
			materiaisIniciais = append(materiaisIniciais, harvested...)
		}
	}

	var materiais []material.Material

	// Intelligent Verification Algorithm & Deduplication
	if len(materiaisIniciais) > 0 && uc.Verifier != nil {
		var urlsToCheck []string
		
		// Deduplicate (harvested might overlap with DB)
		uniqueMap := make(map[string]material.Material)
		for _, m := range materiaisIniciais {
			uID := m.ExternoID
			if uID == "" {
				uID = m.Titulo + ":" + m.Autor
			}
			if _, exists := uniqueMap[uID]; !exists {
				uniqueMap[uID] = m
				if m.PDFURL != "" {
					urlsToCheck = append(urlsToCheck, m.PDFURL)
				}
			}
		}

		// Verify URLs with timeout
		verifyCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		statusMap := uc.Verifier.VerifyBatch(verifyCtx, urlsToCheck)

		for _, m := range uniqueMap {
			if m.PDFURL != "" {
				if alive, exists := statusMap[m.PDFURL]; exists && !alive {
					continue
				}
			}
			materiais = append(materiais, m)
			if limit > 0 && len(materiais) >= limit {
				break
			}
		}
	} else {
		materiais = materiaisIniciais
		if limit > 0 && len(materiais) > limit {
			materiais = materiais[:limit]
		}
	}

	if uc.Cache != nil && len(materiais) > 0 {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
