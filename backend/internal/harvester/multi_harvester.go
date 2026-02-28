package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/pkg/metadata"
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type MultiSourceHarvester struct {
	scielo *SciELOHarvester
	capes  *CAPESHarvester
	ieee   *IEEEHarvester
	meta   *metadata.MetadataService
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		scielo: NewSciELOHarvester(),
		capes:  NewCAPESHarvester(),
		ieee:   NewIEEEHarvester(),
		meta:   metadata.NewMetadataService(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	// Individual harvester timeout to prevent one slow source from blocking the whole request
	// But we still respect the parent context's lifecycle
	searchCtx, cancel := context.WithTimeout(ctx, 5000*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var allMaterials []material.Material
	seen := make(map[string]bool)

	// Target mapping for cleaner loop
	targetHarvesters := []struct {
		name string
		f    func(context.Context, string, string, int) ([]material.Material, error)
	}{
		{"SciELO", h.scielo.Search},
		{"CAPES", h.capes.Search},
		{"IEEE", h.ieee.Search},
	}

	for _, harv := range targetHarvesters {
		if source != "" && !strings.Contains(strings.ToLower(harv.name), strings.ToLower(source)) {
			continue
		}
		wg.Add(1)
		go func(name string, searchFn func(context.Context, string, string, int) ([]material.Material, error)) {
			defer wg.Done()

			enhancedQuery := query
			if startYear > 0 || endYear > 0 {
				enhancedQuery = fmt.Sprintf("%s year:%d-%d", query, startYear, endYear)
			}

			mats, err := searchFn(searchCtx, enhancedQuery, category, limit)
			if err != nil {
				// We don't return error here to allow other harvesters to complete
				logger.Warn("Harvester failed", zap.String("source", name), zap.Error(err))
				return
			}

			mu.Lock()
			for _, m := range mats {
				id := m.ExternoID
				if id == "" {
					id = m.Titulo
				}
				if !seen[id] {
					seen[id] = true
					// Tenta enriquecer se faltar capa ou descrição
					if m.CapaURL == "" || m.Descricao == "" {
						cover, desc := h.meta.FetchEnrichment(m.Titulo, m.Autor, m.ISBN)
						if m.CapaURL == "" && cover != "" {
							m.CapaURL = cover
						}
						if m.Descricao == "" && desc != "" {
							m.Descricao = desc
						}
					}
					allMaterials = append(allMaterials, m)
				}
			}
			mu.Unlock()
		}(harv.name, harv.f)
	}

	wg.Wait()
	logger.Info("Multi-source search completed", zap.Int("total_results", len(allMaterials)))
	return allMaterials, nil
}
