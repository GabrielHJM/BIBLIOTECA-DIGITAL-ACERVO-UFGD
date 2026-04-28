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
	capes    *CAPESHarvester
	ss       *SemanticScholarHarvester
	arxiv    *ArXivHarvester
	gb       *GoogleBooksHarvester
	ol       *OpenLibraryHarvester
	gut      *GutendexHarvester
	doaj     *DOAJHarvester
	ia       *InternetArchiveHarvester
	cross    *CrossrefHarvester
	epmc     *EuropePMCHarvester
	dblp     *DBLPHarvester
	plos     *PLOSHarvester
	openalex *OpenAlexHarvester
	zenodo   *ZenodoHarvester
	hal      *HALHarvester
	pubmed   *PubMedHarvester
	osf      *OSFHarvester
	base     *BASEHarvester
	core     *COREHarvester
	scielo   *SciELOHarvester
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	GlobalSupervisor.RegisterAPI("GoogleBooks", "https://www.googleapis.com/books/v1/volumes?q=test")
	GlobalSupervisor.RegisterAPI("SemanticScholar", "https://api.semanticscholar.org/graph/v1/paper/search?query=test")
	GlobalSupervisor.RegisterAPI("ArXiv", "http://export.arxiv.org/api/query?search_query=all:test")
	GlobalSupervisor.RegisterAPI("CAPES", "https://api.crossref.org/works?query=test")
	GlobalSupervisor.RegisterAPI("OpenLibrary", "https://openlibrary.org/search.json?q=test")
	GlobalSupervisor.RegisterAPI("Gutendex", "https://gutendex.com/books/?search=test")
	GlobalSupervisor.RegisterAPI("DOAJ", "https://doaj.org/api/search/articles/test")
	GlobalSupervisor.RegisterAPI("InternetArchive", "https://archive.org/advancedsearch.php?q=test")
	GlobalSupervisor.RegisterAPI("Crossref", "https://api.crossref.org/works?query=test")
	GlobalSupervisor.RegisterAPI("EuropePMC", "https://www.ebi.ac.uk/europepmc/webservices/rest/search?query=test")
	GlobalSupervisor.RegisterAPI("DBLP", "https://dblp.org/search/publ/api?q=test")
	GlobalSupervisor.RegisterAPI("PLOS", "https://api.plos.org/search?q=test")
	GlobalSupervisor.RegisterAPI("OpenAlex", "https://api.openalex.org/works?search=test")
	GlobalSupervisor.RegisterAPI("Zenodo", "https://zenodo.org/api/records?q=test")
	GlobalSupervisor.RegisterAPI("HAL", "https://api.archives-ouvertes.fr/search/?q=test")
	GlobalSupervisor.RegisterAPI("PubMed", "https://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?db=pmc&term=test")
	GlobalSupervisor.RegisterAPI("OSF", "https://api.osf.io/v2/nodes/?filter[title]=test")
	GlobalSupervisor.RegisterAPI("BASE", "https://api.base-search.net/cgi-bin/BaseHttpSearch/v1/")
	GlobalSupervisor.RegisterAPI("CORE", "https://api.core.ac.uk/v3/search/works")
	GlobalSupervisor.RegisterAPI("SciELO", "https://search.scielo.org/")

	return &MultiSourceHarvester{
		capes:    NewCAPESHarvester(),
		ss:       NewSemanticScholarHarvester(),
		arxiv:    NewArXivHarvester(),
		gb:       NewGoogleBooksHarvester(),
		ol:       NewOpenLibraryHarvester(),
		gut:      NewGutendexHarvester(),
		doaj:     NewDOAJHarvester(),
		ia:       NewInternetArchiveHarvester(),
		cross:    NewCrossrefHarvester(),
		epmc:     NewEuropePMCHarvester(),
		dblp:     NewDBLPHarvester(),
		plos:     NewPLOSHarvester(),
		openalex: NewOpenAlexHarvester(),
		zenodo:   NewZenodoHarvester(),
		hal:      NewHALHarvester(),
		pubmed:   NewPubMedHarvester(),
		osf:      NewOSFHarvester(),
		base:     NewBASEHarvester(),
		core:     NewCOREHarvester(),
		scielo:   NewSciELOHarvester(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int, offset int) ([]material.Material, error) {
	refinedQuery := query
	lowercaseQ := strings.ToLower(query)
	lowercaseC := strings.ToLower(category)

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
 
	pagesToSweep := 2
	if limit > 20 {
		pagesToSweep = 4
	}

	startPage := (offset / 40) + 1

	tasks := []func(ctx context.Context){
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("GoogleBooks") { return }
			mats, err := h.gb.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("SemanticScholar") { return }
			mats, err := h.ss.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("ArXiv") { return }
			mats, err := h.arxiv.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("CAPES") { return }
			mats, err := h.capes.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("OpenLibrary") { return }
			for p := startPage; p < startPage+pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					mats, err := h.ol.Search(c, refinedQuery, category, page, limit/2)
					if err == nil { resultsChan <- mats }
				}(p)
			}
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("Gutendex") { return }
			for p := startPage; p < startPage+pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					mats, err := h.gut.Search(c, refinedQuery, category, page)
					if err == nil { resultsChan <- mats }
				}(p)
			}
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("DOAJ") { return }
			mats, err := h.doaj.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("InternetArchive") { return }
			mats, err := h.ia.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("Crossref") { return }
			mats, err := h.cross.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("EuropePMC") { return }
			mats, err := h.epmc.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("DBLP") { return }
			mats, err := h.dblp.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("PLOS") { return }
			mats, err := h.plos.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		// The 8 new APIs
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("OpenAlex") { return }
			mats, err := h.openalex.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("Zenodo") { return }
			mats, err := h.zenodo.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("HAL") { return }
			mats, err := h.hal.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("PubMed") { return }
			mats, err := h.pubmed.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("OSF") { return }
			mats, err := h.osf.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("BASE") { return }
			mats, err := h.base.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("CORE") { return }
			mats, err := h.core.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("SciELO") { return }
			mats, err := h.scielo.Search(c, refinedQuery, category, limit, offset)
			if err == nil { resultsChan <- mats }
		},
	}

	searchCtx, cancel := context.WithTimeout(ctx, 15*time.Second) // Increased timeout slightly for more APIs
	defer cancel()

	for _, task := range tasks {
		wg.Add(1)
		go func(t func(context.Context)) {
			defer wg.Done()
			t(searchCtx)
		}(task)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for mats := range resultsChan {
		allMaterials = append(allMaterials, mats...)
	}

	uniqueMaterials := make([]material.Material, 0)
	seen := make(map[string]bool)

	for _, m := range allMaterials {
		if m.Titulo == "" || m.PDFURL == "" || !m.Disponivel {
			continue
		}

		lowerURL := strings.ToLower(m.PDFURL)
		isValidReader := strings.HasSuffix(lowerURL, ".pdf") || 
			strings.Contains(lowerURL, "reader") || 
			strings.Contains(lowerURL, "view") || 
			strings.Contains(lowerURL, "archive.org") ||
			strings.Contains(lowerURL, "books.google") ||
			strings.Contains(lowerURL, "gutendex.com") ||
			strings.Contains(lowerURL, "arxiv.org") ||
			strings.Contains(lowerURL, "crossref.org") ||
			strings.Contains(lowerURL, "europepmc.org") ||
			strings.Contains(lowerURL, "dblp.org") ||
			strings.Contains(lowerURL, "plos.org") ||
			strings.Contains(lowerURL, "doi.org") ||
			strings.Contains(lowerURL, "openalex.org") ||
			strings.Contains(lowerURL, "zenodo.org") ||
			strings.Contains(lowerURL, "archives-ouvertes.fr") ||
			strings.Contains(lowerURL, "nih.gov") ||
			strings.Contains(lowerURL, "osf.io")

		if !isValidReader {
			continue
		}

		sig := strings.ToLower(m.Titulo + ":" + m.Autor)
		if m.ExternoID != "" {
			sig = m.ExternoID
		}

		if !seen[sig] {
			seen[sig] = true
			uniqueMaterials = append(uniqueMaterials, m)
		}
	}

	logger.Info("HighCapacity search completed with Supervisor checks", 
		zap.Int("total_harvested", len(allMaterials)), 
		zap.Int("unique_high_quality", len(uniqueMaterials)),
		zap.String("query", query))

	return uniqueMaterials, nil
}
