package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/cache"
	"biblioteca-digital-api/internal/pkg/utils"
	"context"
	"fmt"
	"sort"
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

func (uc *PesquisarMaterialUseCase) Execute(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sortParam string) ([]material.Material, error) {
	cacheKey := fmt.Sprintf("search:%s:%s:%s:%d:%d:%d:%d:%s", termo, categoria, fonte, anoInicio, anoFim, limit, offset, sortParam)
	if uc.Cache != nil {
		var cached []material.Material
		if found := uc.Cache.Get(cacheKey, &cached); found {
			return cached, nil
		}
	}

	// 1. Local Search (Database)
	// We oversample to compensate for any later filtering
	fetchLimit := limit + (limit / 2)
	materiaisIniciais, localErr := uc.Repo.Pesquisar(ctx, termo, categoria, fonte, anoInicio, anoFim, tags, fetchLimit, offset, sortParam)
	if localErr != nil {
		return nil, localErr
	}

	// 2. High Density Harvesting (Trigger if local database is too sparse)
	// We define "sparse" as having less than half of the requested limit OR if the user specifically asked for a high limit
	isSparse := len(materiaisIniciais) < limit || limit > 50 
	hasQuery := termo != "" || categoria != ""

	if hasQuery && isSparse && uc.Harvester != nil && offset == 0 {
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

	if sortParam == "relevancia" || sortParam == "" {
		type ScoredMaterial struct {
			m     material.Material
			score float64
		}

		scored := make([]ScoredMaterial, 0, len(materiais))
		normalizedTermo := strings.ToLower(utils.RemoveAccents(termo))

		for _, m := range materiais {
			score := 0.0
			
			score += float64(m.Relevancia) * 1.0
			
			mTitleLow := strings.ToLower(utils.RemoveAccents(m.Titulo))
			mAuthorLow := strings.ToLower(utils.RemoveAccents(m.Autor))
			
			if normalizedTermo != "" {
				if strings.HasPrefix(mTitleLow, normalizedTermo) { score += 100 }
				if strings.Contains(mTitleLow, normalizedTermo) { score += 50 }
				if strings.Contains(mAuthorLow, normalizedTermo) { score += 30 }
			}

			if m.CapaURL != "" { score += 20 }
			if m.Descricao != "" && len(m.Descricao) > 50 { score += 10 }
			if m.MediaNota > 4.0 { score += 15 }
			
			yearDiff := time.Now().Year() - m.AnoPublicacao
			if yearDiff < 0 { yearDiff = 0 }
			
			// Granular Freshness Boost: Newer books get progressively more points
			if yearDiff < 10 {
				score += float64(30 - (yearDiff * 3)) // 0=30pts, 1=27pts, 2=24pts...
			} else if yearDiff < 20 {
				score += 5
			}

			scored = append(scored, ScoredMaterial{m, score})
		}

		sort.SliceStable(scored, func(i, j int) bool {
			return scored[i].score > scored[j].score
		})

		materiais = make([]material.Material, 0, len(scored))
		for _, sm := range scored {
			materiais = append(materiais, sm.m)
		}
	} else {
		// Respect original sort constraints
		sort.SliceStable(materiais, func(i, j int) bool {
			switch sortParam {
			case "az":
				a := strings.ToLower(utils.RemoveAccents(materiais[i].Titulo))
				b := strings.ToLower(utils.RemoveAccents(materiais[j].Titulo))
				if a == b { return materiais[i].ID > materiais[j].ID }
				return a < b
			case "za":
				a := strings.ToLower(utils.RemoveAccents(materiais[i].Titulo))
				b := strings.ToLower(utils.RemoveAccents(materiais[j].Titulo))
				if a == b { return materiais[i].ID > materiais[j].ID }
				return a > b
			case "recent":
				if materiais[i].AnoPublicacao == materiais[j].AnoPublicacao {
					return materiais[i].ID > materiais[j].ID
				}
				return materiais[i].AnoPublicacao > materiais[j].AnoPublicacao
			case "oldest":
				if materiais[i].AnoPublicacao == materiais[j].AnoPublicacao {
					return materiais[i].ID > materiais[j].ID
				}
				return materiais[i].AnoPublicacao < materiais[j].AnoPublicacao
			}
			return false
		})
	}

	finalMaterials := materiais
	if limit > 0 && len(finalMaterials) > limit {
		finalMaterials = finalMaterials[:limit]
	}

	// 5. Caching
	if uc.Cache != nil && len(finalMaterials) > 0 {
		uc.Cache.Set(cacheKey, finalMaterials, 10*time.Minute)
	}

	return finalMaterials, nil
}
