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

	"go.uber.org/zap"
)

// DOAJ (Directory of Open Access Journals) Harvester
type DOAJHarvester struct {
	BaseURL string
}

func NewDOAJHarvester() *DOAJHarvester {
	return &DOAJHarvester{
		BaseURL: "https://doaj.org/api/v2/search/articles",
	}
}

type DOAJResponse struct {
	Results []struct {
		Bibjson struct {
			Title    string `json:"title"`
			Abstract string `json:"abstract"`
			Year     string `json:"year"`
			Authors  []struct {
				Name string `json:"name"`
			} `json:"author"`
			Link []struct {
				URL  string `json:"url"`
				Type string `json:"type"`
			} `json:"link"`
			Subject []struct {
				Term string `json:"term"`
			} `json:"subject"`
		} `json:"bibjson"`
		ID string `json:"id"`
	} `json:"results"`
}

func (h *DOAJHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "academic"
	}

	searchURL := fmt.Sprintf("%s/%s?pageSize=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("DOAJ harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doaj api error: %s", resp.Status)
	}

	var data DOAJResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Results {
		bib := item.Bibjson
		if bib.Title == "" {
			continue
		}

		var authors []string
		for _, a := range bib.Authors {
			authors = append(authors, a.Name)
		}

		year := 0
		fmt.Sscanf(bib.Year, "%d", &year)

		var pdfURL string
		for _, link := range bib.Link {
			// DOAJ often provides fulltext links and explicit PDF types
			if link.Type == "fulltext" && strings.Contains(strings.ToLower(link.URL), ".pdf") {
				pdfURL = link.URL
				break
			}
		}

		// If no direct .pdf found, we try to take the first fulltext link that might be a PDF redirect
		if pdfURL == "" {
			for _, link := range bib.Link {
				if link.Type == "fulltext" {
					pdfURL = link.URL
					break
				}
			}
		}

		if pdfURL == "" {
			continue
		}

		catName := category
		if catName == "" {
			if len(bib.Subject) > 0 {
				catName = strings.ToUpper(bib.Subject[0].Term)
			} else {
				catName = "ACADÊMICO"
			}
		}

		m := material.Material{
			Titulo:        bib.Title,
			Autor:         strings.Join(authors, ", "),
			Descricao:     bib.Abstract,
			AnoPublicacao: year,
			Fonte:         "DOAJ",
			Categoria:     catName,
			ExternoID:     item.ID,
			PDFURL:        pdfURL,
			Disponivel:    true,
			CapaURL:       GetCoverFromGoogleBooks(bib.Title, ""),
		}

		materials = append(materials, m)
	}

	logger.Info("DOAJ harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
