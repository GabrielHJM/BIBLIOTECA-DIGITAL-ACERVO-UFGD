package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"go.uber.org/zap"
)

// ==========================================
// OpenAlex Harvester
// ==========================================
type OpenAlexHarvester struct{ client *http.Client }

func NewOpenAlexHarvester() *OpenAlexHarvester {
	return &OpenAlexHarvester{client: &http.Client{Timeout: 10 * time.Second}}
}

func (h *OpenAlexHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	apiURL := fmt.Sprintf("https://api.openalex.org/works?search=%s&per-page=%d&page=%d", url.QueryEscape(query), limit, page)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	
	// OpenAlex politeness
	req.Header.Set("User-Agent", "mailto:contato@bibliotecadigital.com")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openalex bad status: %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			ID               string `json:"id"`
			Title            string `json:"title"`
			PublicationYear  int    `json:"publication_year"`
			Authorships      []struct {
				Author struct {
					DisplayName string `json:"display_name"`
				} `json:"author"`
			} `json:"authorships"`
			OpenAccess struct {
				IsOA  bool   `json:"is_oa"`
				OAUrl string `json:"oa_url"`
			} `json:"open_access"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var mats []material.Material
	for _, item := range data.Results {
		if item.Title == "" || !item.OpenAccess.IsOA || item.OpenAccess.OAUrl == "" {
			continue
		}

		pdfURL := item.OpenAccess.OAUrl
		if !strings.HasSuffix(strings.ToLower(pdfURL), ".pdf") {
			// Muitas vezes o OpenAlex retorna páginas HTML no oa_url. Vamos pular se não parecer PDF.
			if !strings.Contains(strings.ToLower(pdfURL), "pdf") {
				continue
			}
		}

		var authors []string
		for _, a := range item.Authorships {
			if a.Author.DisplayName != "" {
				authors = append(authors, a.Author.DisplayName)
			}
		}
		authorStr := strings.Join(authors, ", ")

		mats = append(mats, material.Material{
			Titulo:        item.Title,
			Autor:         authorStr,
			Descricao:     "Material acadêmico extraído de OpenAlex.",
			AnoPublicacao: item.PublicationYear,
			Categoria:     category,
			PDFURL:        pdfURL,
			Fonte:         "OpenAlex",
			Disponivel:    true,
			Dificuldade:   3,
			XP:            25,
			Relevancia:    15,
		})
	}
	logger.Info("OpenAlex search completed", zap.Int("results", len(mats)))
	return mats, nil
}

// ==========================================
// Zenodo Harvester
// ==========================================
type ZenodoHarvester struct{ client *http.Client }

func NewZenodoHarvester() *ZenodoHarvester {
	return &ZenodoHarvester{client: &http.Client{Timeout: 10 * time.Second}}
}

func (h *ZenodoHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	page := (offset / limit) + 1
	if page < 1 {
		page = 1
	}
	apiURL := fmt.Sprintf("https://zenodo.org/api/records?q=%s&size=%d&page=%d&access_right=open", url.QueryEscape(query), limit, page)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("zenodo bad status: %d", resp.StatusCode)
	}

	var data struct {
		Hits struct {
			Hits []struct {
				Metadata struct {
					Title       string `json:"title"`
					Description string `json:"description"`
					PublicationDate string `json:"publication_date"`
					Creators []struct {
						Name string `json:"name"`
					} `json:"creators"`
				} `json:"metadata"`
				Files []struct {
					Links struct {
						Self string `json:"self"`
					} `json:"links"`
					Type string `json:"type"`
				} `json:"files"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var mats []material.Material
	for _, item := range data.Hits.Hits {
		if item.Metadata.Title == "" {
			continue
		}

		var pdfURL string
		for _, f := range item.Files {
			if strings.ToLower(f.Type) == "pdf" || strings.HasSuffix(strings.ToLower(f.Links.Self), ".pdf") {
				pdfURL = f.Links.Self
				break
			}
		}

		if pdfURL == "" {
			continue
		}

		var authors []string
		for _, c := range item.Metadata.Creators {
			authors = append(authors, c.Name)
		}

		year := 0
		if len(item.Metadata.PublicationDate) >= 4 {
			fmt.Sscanf(item.Metadata.PublicationDate[:4], "%d", &year)
		}

		// Limpar HTML do description
		desc := strings.ReplaceAll(item.Metadata.Description, "<p>", "")
		desc = strings.ReplaceAll(desc, "</p>", "")

		mats = append(mats, material.Material{
			Titulo:        item.Metadata.Title,
			Autor:         strings.Join(authors, ", "),
			Descricao:     desc,
			AnoPublicacao: year,
			Categoria:     category,
			PDFURL:        pdfURL,
			Fonte:         "Zenodo",
			Disponivel:    true,
			Dificuldade:   3,
			XP:            25,
			Relevancia:    15,
		})
	}
	logger.Info("Zenodo search completed", zap.Int("results", len(mats)))
	return mats, nil
}

// ==========================================
// HAL Harvester (Hyper Articles en Ligne)
// ==========================================
type HALHarvester struct{ client *http.Client }

func NewHALHarvester() *HALHarvester {
	return &HALHarvester{client: &http.Client{Timeout: 10 * time.Second}}
}

func (h *HALHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	apiURL := fmt.Sprintf("https://api.archives-ouvertes.fr/search/?q=%s&rows=%d&start=%d&wt=json", url.QueryEscape(query), limit, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("hal bad status: %d", resp.StatusCode)
	}

	var data struct {
		Response struct {
			Docs []struct {
				TitleS         []string `json:"title_s"`
				AuthFullNameS  []string `json:"authFullName_s"`
				ProducedDateYI int      `json:"producedDateY_i"`
				FileMainS      string   `json:"fileMain_s"`
				AbstractS      []string `json:"abstract_s"`
			} `json:"docs"`
		} `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var mats []material.Material
	for _, doc := range data.Response.Docs {
		if len(doc.TitleS) == 0 || doc.FileMainS == "" {
			continue
		}

		desc := ""
		if len(doc.AbstractS) > 0 {
			desc = doc.AbstractS[0]
		}

		mats = append(mats, material.Material{
			Titulo:        doc.TitleS[0],
			Autor:         strings.Join(doc.AuthFullNameS, ", "),
			Descricao:     desc,
			AnoPublicacao: doc.ProducedDateYI,
			Categoria:     category,
			PDFURL:        doc.FileMainS,
			Fonte:         "HAL",
			Disponivel:    true,
			Dificuldade:   4,
			XP:            30,
			Relevancia:    18,
		})
	}
	logger.Info("HAL search completed", zap.Int("results", len(mats)))
	return mats, nil
}

// ==========================================
// SciELO Harvester (Via Crossref)
// ==========================================
type SciELOHarvester struct{ client *http.Client }

func NewSciELOHarvester() *SciELOHarvester {
	return &SciELOHarvester{client: &http.Client{Timeout: 10 * time.Second}}
}

func (h *SciELOHarvester) Search(ctx context.Context, query, category string, limit, offset int) ([]material.Material, error) {
	// SciELO uses DOI prefix 10.1590
	apiURL := fmt.Sprintf("https://api.crossref.org/works?query=%s&filter=prefix:10.1590,has-full-text:true&rows=%d&offset=%d", url.QueryEscape(query), limit, offset)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scielo bad status: %d", resp.StatusCode)
	}

	var data struct {
		Message struct {
			Items []struct {
				Title []string `json:"title"`
				Author []struct {
					Given  string `json:"given"`
					Family string `json:"family"`
				} `json:"author"`
				Issued struct {
					DateParts [][]int `json:"date-parts"`
				} `json:"issued"`
				Abstract string `json:"abstract"`
				Link []struct {
					URL         string `json:"URL"`
					ContentType string `json:"content-type"`
				} `json:"link"`
			} `json:"items"`
		} `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var mats []material.Material
	for _, item := range data.Message.Items {
		if len(item.Title) == 0 {
			continue
		}

		var pdfURL string
		for _, l := range item.Link {
			if l.ContentType == "application/pdf" {
				pdfURL = l.URL
				break
			}
		}

		if pdfURL == "" {
			continue
		}

		var authors []string
		for _, a := range item.Author {
			authors = append(authors, fmt.Sprintf("%s %s", a.Given, a.Family))
		}

		year := 0
		if len(item.Issued.DateParts) > 0 && len(item.Issued.DateParts[0]) > 0 {
			year = item.Issued.DateParts[0][0]
		}

		mats = append(mats, material.Material{
			Titulo:        item.Title[0],
			Autor:         strings.Join(authors, ", "),
			Descricao:     item.Abstract,
			AnoPublicacao: year,
			Categoria:     category,
			PDFURL:        pdfURL,
			Fonte:         "SciELO",
			Disponivel:    true,
			Dificuldade:   4,
			XP:            35, // SciELO local content gets a bit more XP!
			Relevancia:    20,
		})
	}
	logger.Info("SciELO search completed", zap.Int("results", len(mats)))
	return mats, nil
}
