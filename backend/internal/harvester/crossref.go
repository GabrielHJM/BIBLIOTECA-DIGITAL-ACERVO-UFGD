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

type CrossrefNewResponse struct {
	Message struct {
		Items []struct {
			DOI    string   `json:"DOI"`
			Title  []string `json:"title"`
			Author []struct {
				Given  string `json:"given"`
				Family string `json:"family"`
			} `json:"author"`
			Link []struct {
				URL                 string `json:"URL"`
				ContentType         string `json:"content-type"`
				IntendedApplication string `json:"intended-application"`
			} `json:"link"`
			Issued struct {
				DateParts [][]int `json:"date-parts"`
			} `json:"issued"`
			Abstract string `json:"abstract"`
		} `json:"items"`
	} `json:"message"`
}

type CrossrefHarvester struct {
	BaseURL string
}

func NewCrossrefHarvester() *CrossrefHarvester {
	return &CrossrefHarvester{
		BaseURL: "https://api.crossref.org/works",
	}
}

func (h *CrossrefHarvester) Search(ctx context.Context, query string, category string, limit int, offset int) ([]material.Material, error) {
	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	searchURL := fmt.Sprintf("%s?query=%s&select=DOI,title,author,abstract,issued,link&rows=%d&offset=%d", h.BaseURL, url.QueryEscape(searchTerm), limit, offset)

	req, errReq := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if errReq != nil {
		return nil, errReq
	}
	
	// Polite pool (optional, but good practice for Crossref: Mailto header)
	req.Header.Set("User-Agent", "AcervusBot/1.0 (mailto:gabriel@example.com)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crossref api error: %s", resp.Status)
	}

	var data CrossrefNewResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Message.Items {
		if len(item.Title) == 0 {
			continue
		}
		title := item.Title[0]

		var pdfURL string
		for _, link := range item.Link {
			if link.ContentType == "application/pdf" {
				pdfURL = link.URL
				break
			}
		}

		if pdfURL == "" {
			continue // Pula se não tiver PDF
		}

		var autores []string
		for _, a := range item.Author {
			autores = append(autores, fmt.Sprintf("%s %s", a.Given, a.Family))
		}
		autor := "Desconhecido"
		if len(autores) > 0 {
			autor = strings.Join(autores, ", ")
		}

		year := 0
		if len(item.Issued.DateParts) > 0 && len(item.Issued.DateParts[0]) > 0 {
			year = item.Issued.DateParts[0][0]
		}

		m := material.Material{
			Titulo:        title,
			Autor:         autor,
			Descricao:     strings.TrimSpace(item.Abstract),
			AnoPublicacao: year,
			Fonte:         "Crossref",
			Categoria:     category,
			ExternoID:     item.DOI,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   4, // Geralmente papers
			XP:            30,
			Relevancia:    15,
		}

		materials = append(materials, m)
	}

	logger.Info("Crossref harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
