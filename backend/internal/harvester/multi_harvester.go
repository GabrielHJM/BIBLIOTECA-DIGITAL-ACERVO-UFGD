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
	capes *CAPESHarvester
	ss    *SemanticScholarHarvester
	arxiv *ArXivHarvester
	gb    *GoogleBooksHarvester
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		capes: NewCAPESHarvester(),
		ss:    NewSemanticScholarHarvester(),
		arxiv: NewArXivHarvester(),
		gb:    NewGoogleBooksHarvester(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	// Refine query for English-based academic databases from PT-BR targets
	refinedQuery := query
	lowercaseQ := strings.ToLower(query)
	lowercaseC := strings.ToLower(category)

	if lowercaseQ == "tecnologia" || lowercaseC == "tecnologia" {
		refinedQuery = "computer science OR software OR technology"
	} else if lowercaseQ == "saúde" || lowercaseC == "saúde" {
		refinedQuery = "health OR medicine OR biology"
	} else if lowercaseQ == "ciências" || lowercaseC == "ciências" {
		refinedQuery = "science OR physics OR chemistry"
	} else if lowercaseQ == "matemática" || lowercaseC == "matemática" {
		refinedQuery = "mathematics OR algebra OR geometry"
	} else if lowercaseQ == "história" || lowercaseC == "história" {
		refinedQuery = "history OR archaeology OR humanity"
	} else if lowercaseQ == "educação" || lowercaseC == "educação" {
		refinedQuery = "education OR pedagogical OR teaching"
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
			mats, err := h.ss.Search(ctx, refinedQuery, category, limit)
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
			mats, err := h.gb.Search(ctx, refinedQuery, category, limit)
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
