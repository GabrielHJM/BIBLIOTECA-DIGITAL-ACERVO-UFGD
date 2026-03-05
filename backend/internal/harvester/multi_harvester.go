package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type MultiSourceHarvester struct {
	capes  *CAPESHarvester
	gbooks *GoogleBooksHarvester
	arxiv  *ArXivHarvester
	doaj   *DOAJHarvester
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		capes:  NewCAPESHarvester(),
		gbooks: NewGoogleBooksHarvester(),
		arxiv:  NewArXivHarvester(),
		doaj:   NewDOAJHarvester(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	// Refine query for technology if it's broad to bring more modern results
	refinedQuery := query
	if strings.ToLower(query) == "tecnologia" || strings.ToLower(category) == "tecnologia" {
		refinedQuery = "tecnologia \"artificial intelligence\" OR \"software engineering\" OR \"cybersecurity\""
	}

	var allMaterials []material.Material
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Harvesters to run
	searchFuncs := []func(){
		func() {
			mats, err := h.capes.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			}
		},
		func() {
			mats, err := h.gbooks.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			}
		},
		func() {
			mats, err := h.arxiv.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			}
		},
		func() {
			mats, err := h.doaj.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			}
		},
	}

	for _, fn := range searchFuncs {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(fn)
	}

	wg.Wait()

	// Deduplicate
	var uniqueMaterials []material.Material
	seen := make(map[string]bool)

	for _, m := range allMaterials {
		id := m.ExternoID
		if id == "" {
			id = m.Titulo + ":" + m.Autor
		}
		if !seen[id] {
			seen[id] = true
			uniqueMaterials = append(uniqueMaterials, m)
		}
	}

	logger.Info("MultiSource academic search completed", zap.Int("total_results", len(uniqueMaterials)))
	return uniqueMaterials, nil
}
