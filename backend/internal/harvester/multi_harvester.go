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

	idx := 0
	if limit > 0 {
		idx = (offset / limit)
	} else {
		idx = offset
	}

	if lowercaseQ == "tecnologia" || lowercaseC == "tecnologia" || lowercaseC == "tecnologia brasil" {
		subs := []string{"tecnologia computação", "engenharia de software", "inteligência artificial", "redes de computadores", "segurança da informação", "programação", "banco de dados"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseQ == "saúde" || lowercaseC == "saúde" || lowercaseC == "saúde pública brasil" {
		subs := []string{"saúde medicina", "enfermagem", "saúde pública", "epidemiologia", "fisioterapia", "nutrição"}
		refinedQuery = subs[idx % len(subs)]
	} else if strings.Contains(lowercaseQ, "odontolog") || strings.Contains(lowercaseC, "odontolog") {
		subs := []string{"odontologia saúde", "odontologia clínica", "cirurgia odontológica", "ortodontia", "periodontia"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseQ == "ciências" || lowercaseC == "ciências" || lowercaseC == "science" {
		subs := []string{"ciências pesquisa", "física", "química", "biologia", "astronomia", "ciências da natureza"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseQ == "matemática" || lowercaseC == "matemática" || lowercaseC == "mathematics" {
		subs := []string{"matemática cálculo", "álgebra", "geometria", "estatística", "matemática aplicada", "probabilidade"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseC == "história" {
		subs := []string{"história mundial", "história do brasil", "história antiga", "idade média", "história contemporânea"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseC == "educação" {
		subs := []string{"educação pedagogia", "didática", "educação infantil", "ensino superior", "educação inclusiva"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseC == "jurídico" || lowercaseC == "direito brasileiro" {
		subs := []string{"direito constitucional", "direito penal", "direito civil", "direito do trabalho", "processo penal", "direito administrativo"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseC == "literatura brasileira" {
		subs := []string{"literatura brasileira", "poesia brasileira", "romance brasileiro", "contos brasileiros", "modernismo brasileiro"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseC == "contabilidade" {
		subs := []string{"contabilidade geral", "contabilidade financeira", "auditoria", "contabilidade pública", "contabilidade de custos"}
		refinedQuery = subs[idx % len(subs)]
	} else if lowercaseQ == "" {
		subjects := []string{
			"tecnologia", "história", "medicina", "literatura", "filosofia", 
			"matemática", "física", "biologia", "direito", "economia", 
			"sociologia", "educação", "geografia", "engenharia", "psicologia", 
			"arte", "música", "astronomia", "química", "arquitetura", "ciência",
		}
		refinedQuery = subjects[idx % len(subjects)]
	}

	var allMaterials []material.Material
	resultsChan := make(chan []material.Material, 100)
	var wg sync.WaitGroup
 
	pagesToSweep := 2
	if limit > 20 {
		pagesToSweep = 4
	}

	startPage := (offset / 40) + 1

	executeWithBackoff := func(c context.Context, apiName string, fetch func() ([]material.Material, error)) {
		if !GlobalSupervisor.IsOnline(apiName) { return }
		
		retries := 3
		backoff := 500 * time.Millisecond
		
		for i := 0; i < retries; i++ {
			if c.Err() != nil { return }
			mats, err := fetch()
			if err == nil {
				if len(mats) > 0 {
					resultsChan <- mats
				}
				return
			}
			select {
			case <-time.After(backoff):
				backoff *= 2
			case <-c.Done():
				return
			}
		}
	}

	tasks := []func(ctx context.Context){
		func(c context.Context) {
			executeWithBackoff(c, "GoogleBooks", func() ([]material.Material, error) { return h.gb.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "SemanticScholar", func() ([]material.Material, error) { return h.ss.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "ArXiv", func() ([]material.Material, error) { return h.arxiv.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "CAPES", func() ([]material.Material, error) { return h.capes.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("OpenLibrary") { return }
			for p := startPage; p < startPage+pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					executeWithBackoff(c, "OpenLibrary", func() ([]material.Material, error) { return h.ol.Search(c, refinedQuery, category, page, limit/2) })
				}(p)
			}
		},
		func(c context.Context) {
			if !GlobalSupervisor.IsOnline("Gutendex") { return }
			for p := startPage; p < startPage+pagesToSweep; p++ {
				wg.Add(1)
				go func(page int) {
					defer wg.Done()
					executeWithBackoff(c, "Gutendex", func() ([]material.Material, error) { return h.gut.Search(c, refinedQuery, category, page) })
				}(p)
			}
		},
		func(c context.Context) {
			executeWithBackoff(c, "DOAJ", func() ([]material.Material, error) { return h.doaj.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "InternetArchive", func() ([]material.Material, error) { return h.ia.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "Crossref", func() ([]material.Material, error) { return h.cross.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "EuropePMC", func() ([]material.Material, error) { return h.epmc.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "DBLP", func() ([]material.Material, error) { return h.dblp.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "PLOS", func() ([]material.Material, error) { return h.plos.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "OpenAlex", func() ([]material.Material, error) { return h.openalex.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "Zenodo", func() ([]material.Material, error) { return h.zenodo.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "HAL", func() ([]material.Material, error) { return h.hal.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "PubMed", func() ([]material.Material, error) { return h.pubmed.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "OSF", func() ([]material.Material, error) { return h.osf.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "BASE", func() ([]material.Material, error) { return h.base.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "CORE", func() ([]material.Material, error) { return h.core.Search(c, refinedQuery, category, limit, offset) })
		},
		func(c context.Context) {
			executeWithBackoff(c, "SciELO", func() ([]material.Material, error) { return h.scielo.Search(c, refinedQuery, category, limit, offset) })
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
