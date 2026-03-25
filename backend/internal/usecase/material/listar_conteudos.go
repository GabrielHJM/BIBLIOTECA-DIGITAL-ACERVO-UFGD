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

	var materiais []material.Material

	// Intelligent Verification Algorithm
	if len(materiaisIniciais) > 0 && uc.Verifier != nil {
		var urlsToCheck []string
		for _, m := range materiaisIniciais {
			if m.PDFURL != "" {
				urlsToCheck = append(urlsToCheck, m.PDFURL)
			}
		}

		verifyCtx, cancel := context.WithTimeout(ctx, 1200*time.Millisecond)
		defer cancel()
		
		statusMap := uc.Verifier.VerifyBatch(verifyCtx, urlsToCheck)

		for _, m := range materiaisIniciais {
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
	}

	// As bucas externas na API síncrona foram removidas de propósito
	// para evitar extrema latência e 429 Too Many Requests.
	// O background worker populará esse acervo automaticamente.

	if uc.Cache != nil {
		uc.Cache.Set(cacheKey, materiais, 15*time.Minute)
	}

	return materiais, nil
}
