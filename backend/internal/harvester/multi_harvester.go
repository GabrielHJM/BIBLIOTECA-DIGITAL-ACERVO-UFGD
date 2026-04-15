package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

type MultiSourceHarvester struct {
	capes *CAPESHarvester
	ss    *SemanticScholarHarvester
	arxiv *ArXivHarvester
	gb    *GoogleBooksHarvester
	ol    *OpenLibraryHarvester
	gut   *GutendexHarvester
	doaj  *DOAJHarvester
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		capes: NewCAPESHarvester(),
		ss:    NewSemanticScholarHarvester(),
		arxiv: NewArXivHarvester(),
		gb:    NewGoogleBooksHarvester(),
		ol:    NewOpenLibraryHarvester(),
		gut:   NewGutendexHarvester(),
		doaj:  NewDOAJHarvester(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	// 1. Query Expansion (Modularized Logic)
	refinedQuery := query
	lowercaseQ := strings.ToLower(query)
	lowercaseC := strings.ToLower(category)

	// Contextual expansion for better academic coverage (Português/Brasil prioritized) - AND logic
	if lowercaseQ == "tecnologia" || lowercaseC == "tecnologia" {
		refinedQuery = "tecnologia AND computação"
	} else if lowercaseQ == "saúde" || lowercaseC == "saúde" {
		refinedQuery = "saúde AND medicina"
	} else if strings.Contains(lowercaseQ, "odontolog") {
		refinedQuery = "odontologia AND saúde"
	} else if lowercaseQ == "ciências" || lowercaseC == "ciências" {
		refinedQuery = "ciências AND pesquisa"
	} else if lowercaseQ == "matemática" || lowercaseC == "matemática" {
		refinedQuery = "matemática AND cálculo"
	} else if strings.Contains(lowercaseQ, "brasil") || strings.Contains(lowercaseC, "brasil") {
		refinedQuery = query + " AND brasil"
	} else if lowercaseQ == "" {
		refinedQuery = "livro AND acadêmico"
	}

	var allMaterials []material.Material
	resultsChan := make(chan []material.Material, 100)
	var wg sync.WaitGroup
 
	// 2. High Capacity Worker Pool Logic
	// We will sweep up to 5 pages for each high-yield source to force more volume
	pagesToSweep := 3
	if limit > 20 {
		pagesToSweep = 5
	}

	// Define tasks for the worker pool
	tasks := []func(ctx context.Context){
		// Google Books Sweep
		func(c context.Context) {
			mats, err := h.gb.Search(c, refinedQuery, category, limit)
			if err == nil { resultsChan <- mats }
		},
		// Semantic Scholar Sweep
		func(c context.Context) {
			mats, err := h.ss.Search(c, refinedQuery, category, limit)
			if err == nil { resultsChan <- mats }
		},
		// ArXiv Sweep
		func(c context.Context) {
			mats, err := h.arxiv.Search(c, refinedQuery, category, limit)
			if err == nil { resultsChan <- mats }
		},
		// CAPES/Crossref
		func(c context.Context) {
			mats, err := h.capes.Search(c, refinedQuery, category, limit)
			if err == nil { resultsChan <- mats }
		},
	// Open Library Sweep (Multiple Pages in Parallel)
		func(c context.Context) {
			for p := 1; p <= pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					mats, err := h.ol.Search(c, refinedQuery, category, page, limit/2)
					if err == nil { resultsChan <- mats }
				}(p)
			}
		},
		// Gutendex Sweep (Multiple Pages in Parallel)
		func(c context.Context) {
			for p := 1; p <= pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					mats, err := h.gut.Search(c, refinedQuery, category, page)
					if err == nil { resultsChan <- mats }
				}(p)
			}
		},
		// DOAJ (Directory of Open Access Journals)
		func(c context.Context) {
			mats, err := h.doaj.Search(c, refinedQuery, category, limit)
			if err == nil { resultsChan <- mats }
		},
	}

	// 3. Execute concurrently with timeout safety
	searchCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	for _, task := range tasks {
		wg.Add(1)
		go func(t func(context.Context)) {
			defer wg.Done()
			t(searchCtx)
		}(task)
	}

	// Close results channel when all workers are done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 4. Collect results
	for mats := range resultsChan {
		allMaterials = append(allMaterials, mats...)
	}

	// 5. Intelligent Deduplication and Filtering
	uniqueMaterials := make([]material.Material, 0)
	seen := make(map[string]bool)

	for _, m := range allMaterials {
		// Clean junk and enforce "Online Reading Only" rule
		if m.Titulo == "" || m.PDFURL == "" {
			continue
		}

		// Final check: URL must look like a PDF or a known online reader
		lowerURL := strings.ToLower(m.PDFURL)
		isValidReader := strings.HasSuffix(lowerURL, ".pdf") || 
			strings.Contains(lowerURL, "reader") || 
			strings.Contains(lowerURL, "view") || 
			strings.Contains(lowerURL, "archive.org/download") ||
			strings.Contains(lowerURL, "books.google") ||
			strings.Contains(lowerURL, "gutendex.com") ||
			strings.Contains(lowerURL, "arxiv.org/pdf")

		if !isValidReader {
			continue
		}

		// Signature for deduplication
		sig := strings.ToLower(m.Titulo + ":" + m.Autor)
		if m.ExternoID != "" {
			sig = m.ExternoID
		}

		if !seen[sig] {
			seen[sig] = true
			uniqueMaterials = append(uniqueMaterials, m)
		}
	}

	logger.Info("HighCapacity search completed", 
		zap.Int("total_harvested", len(allMaterials)), 
		zap.Int("unique_high_quality", len(uniqueMaterials)),
		zap.String("query", query))

	return uniqueMaterials, nil
}
