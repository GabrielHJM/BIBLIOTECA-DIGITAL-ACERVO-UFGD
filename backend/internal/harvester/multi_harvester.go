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
		refinedQuery = "tecnologia OR computer science OR software OR technology OR artificial intelligence OR web development"
	} else if lowercaseQ == "saúde" || lowercaseC == "saúde" {
		refinedQuery = "saúde OR medicina OR health OR medicine OR nursing OR pharmacy OR public health"
	} else if lowercaseQ == "ciências" || lowercaseC == "ciências" {
		refinedQuery = "ciências OR science OR physics OR chemistry OR biology OR astronomy"
	} else if lowercaseQ == "matemática" || lowercaseC == "matemática" {
		refinedQuery = "matemática OR mathematics OR algebra OR calculus OR statistics"
	} else if lowercaseQ == "história" || lowercaseC == "história" {
		refinedQuery = "história OR history OR archaeology OR humanity OR anthropology"
	} else if lowercaseQ == "educação" || lowercaseC == "educação" {
		refinedQuery = "educação OR pedagogical OR teaching OR higher education OR learning"
	} else if lowercaseQ == "jurídico" || lowercaseC == "jurídico" {
		refinedQuery = "direito OR jurídico OR law OR legal OR constitution OR regulation"
	} else if lowercaseQ == "contabilidade" || lowercaseC == "contabilidade" {
		refinedQuery = "contabilidade OR accounting OR finance OR audit OR taxation"
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
			} else {
				logger.Error("CAPES Harvester falhou", zap.Error(err))
			}
		},
		func() {
			mats, err := h.ss.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			} else {
				logger.Error("SemanticScholar Harvester falhou", zap.Error(err))
			}
		},
		func() {
			mats, err := h.arxiv.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			} else {
				logger.Error("ArXiv Harvester falhou", zap.Error(err))
			}
		},
		func() {
			mats, err := h.gb.Search(ctx, refinedQuery, category, limit)
			if err == nil {
				mu.Lock()
				allMaterials = append(allMaterials, mats...)
				mu.Unlock()
			} else {
				logger.Error("GoogleBooks Harvester falhou", zap.Error(err))
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
		// Global Strict PDF Filter
		lowerLink := strings.ToLower(m.PDFURL)
		isPDF := strings.HasSuffix(strings.Split(lowerLink, "?")[0], ".pdf") ||
			strings.Contains(lowerLink, "pdf") ||
			strings.Contains(lowerLink, "download") ||
			strings.Contains(lowerLink, "googleapis.com")

		if !isPDF {
			continue
		}

		id := m.ExternoID
		if id == "" {
			id = m.Titulo + ":" + m.Autor
		}
		if !seen[id] {
			seen[id] = true
			uniqueMaterials = append(uniqueMaterials, m)
		}
	}

	logger.Info("MultiSource academic search completed", zap.Int("total_filtered_results", len(uniqueMaterials)), zap.String("query", query), zap.String("category", category))
	return uniqueMaterials, nil
}
