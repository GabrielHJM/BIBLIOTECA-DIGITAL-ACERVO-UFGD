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
	ia    *InternetArchiveHarvester
	cross *CrossrefHarvester
	epmc  *EuropePMCHarvester
	dblp  *DBLPHarvester
	plos  *PLOSHarvester
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
		ia:    NewInternetArchiveHarvester(),
		cross: NewCrossrefHarvester(),
		epmc:  NewEuropePMCHarvester(),
		dblp:  NewDBLPHarvester(),
		plos:  NewPLOSHarvester(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int, offset int) ([]material.Material, error) {
	// 1. Query Expansion (Modularized Logic)
	refinedQuery := query
	lowercaseQ := strings.ToLower(query)
	lowercaseC := strings.ToLower(category)

	// Contextual expansion for better academic coverage (Português/Brasil prioritized)
	if lowercaseQ == "tecnologia" || lowercaseC == "tecnologia" {
		refinedQuery = "tecnologia computação"
	} else if lowercaseQ == "saúde" || lowercaseC == "saúde" {
		refinedQuery = "saúde medicina"
	} else if strings.Contains(lowercaseQ, "odontolog") {
		refinedQuery = "odontologia saúde"
	} else if lowercaseQ == "ciências" || lowercaseC == "ciências" || lowercaseC == "science" {
		refinedQuery = "ciências pesquisa"
	} else if lowercaseQ == "matemática" || lowercaseC == "matemática" || lowercaseC == "mathematics" {
		refinedQuery = "matemática cálculo"
	} else if lowercaseQ == "" {
		refinedQuery = "livro"
	}

	var allMaterials []material.Material
	resultsChan := make(chan []material.Material, 100)
	var wg sync.WaitGroup
 
	// 2. High Capacity Worker Pool Logic
	// Reduce pages to sweep for significantly faster results while still keeping volume decent
	pagesToSweep := 2
	if limit > 20 {
		pagesToSweep = 4
	}

	// Calculate start page for APIs that use pages instead of offsets
	// Assume an average of 40 items per page for providers like OpenLibrary
	startPage := (offset / 40) + 1

	// Define tasks for the worker pool
	tasks := []func(ctx context.Context){
		// Google Books Sweep
		func(c context.Context) {
			mats, err := h.gb.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// Semantic Scholar Sweep
		func(c context.Context) {
			mats, err := h.ss.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// ArXiv Sweep
		func(c context.Context) {
			mats, err := h.arxiv.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// CAPES/Crossref
		func(c context.Context) {
			mats, err := h.capes.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
	// Open Library Sweep (Multiple Pages in Parallel)
		func(c context.Context) {
			for p := startPage; p < startPage+pagesToSweep; p++ {
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
			for p := startPage; p < startPage+pagesToSweep; p++ {
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
			mats, err := h.doaj.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// Internet Archive Sweep
		func(c context.Context) {
			mats, err := h.ia.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// Crossref Sweep
		func(c context.Context) {
			mats, err := h.cross.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// Europe PMC Sweep
		func(c context.Context) {
			mats, err := h.epmc.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// DBLP Sweep
		func(c context.Context) {
			mats, err := h.dblp.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// PLOS Sweep
		func(c context.Context) {
			mats, err := h.plos.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
	}

	// 3. Execute concurrently with timeout safety (Fast fail for responsive UI)
	searchCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
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
		if m.Titulo == "" || m.PDFURL == "" || !m.Disponivel {
			continue
		}

		// Restrict to proper browser readers/PDFs to prevent redirecting to external homepages
		lowerURL := strings.ToLower(m.PDFURL)
		isValidReader := strings.HasSuffix(lowerURL, ".pdf") || 
			strings.Contains(lowerURL, "reader") || 
			strings.Contains(lowerURL, "view") || 
			strings.Contains(lowerURL, "archive.org/download") ||
			strings.Contains(lowerURL, "archive.org/stream") ||
			strings.Contains(lowerURL, "books.google") ||
			strings.Contains(lowerURL, "gutendex.com") ||
			strings.Contains(lowerURL, "arxiv.org/pdf") ||
			strings.Contains(lowerURL, "crossref.org") ||
			strings.Contains(lowerURL, "europepmc.org") ||
			strings.Contains(lowerURL, "dblp.org") ||
			strings.Contains(lowerURL, "plos.org") ||
			strings.Contains(lowerURL, "doi.org")

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
