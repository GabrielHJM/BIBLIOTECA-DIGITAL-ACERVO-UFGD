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

// Crossref structures for academic metadata
type CrossrefResponse struct {
	Message struct {
		Items []struct {
			Title  []string `json:"title"`
			Author []struct {
				Given  string `json:"given"`
				Family string `json:"family"`
			} `json:"author"`
			Abstract string `json:"abstract"`
			Created  struct {
				DateParts [][]int `json:"date-parts"`
			} `json:"created"`
			DOI  string `json:"DOI"`
			URL  string `json:"URL"`
			Type string `json:"type"`
		} `json:"items"`
	} `json:"message"`
}

type CAPESHarvester struct {
	BaseURL string
}

func NewCAPESHarvester() *CAPESHarvester {
	return &CAPESHarvester{
		BaseURL: "https://api.crossref.org/works",
	}
}

func (h *CAPESHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?query=%s&rows=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("CAPES/Crossref harvester: request failed", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crossref api error: %s", resp.Status)
	}

	var data CrossrefResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Message.Items {
		if len(item.Title) == 0 {
			continue
		}

		var authors []string
		for _, a := range item.Author {
			authors = append(authors, fmt.Sprintf("%s %s", a.Given, a.Family))
		}

		year := 0
		if len(item.Created.DateParts) > 0 && len(item.Created.DateParts[0]) > 0 {
			year = item.Created.DateParts[0][0]
		}

		catName := category
		if catName == "" {
			catName = "Artigo Periódico"
		}

		m := material.Material{
			Titulo:        item.Title[0],
			Autor:         strings.Join(authors, ", "),
			Descricao:     item.Abstract,
			AnoPublicacao: year,
			Fonte:         "CAPES/Crossref",
			Categoria:     catName,
			ExternoID:     item.DOI,
			PDFURL:        item.URL,
			Disponivel:    true,
		}

		materials = append(materials, m)
	}

	logger.Info("CAPES/Crossref harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
