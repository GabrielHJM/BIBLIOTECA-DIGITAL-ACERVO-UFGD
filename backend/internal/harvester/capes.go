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
	"sync"

	"go.uber.org/zap"
)

// Crossref structures for academic metadata
type CrossrefResponse struct {
	Message struct {
		Items []struct {
			Title   []string `json:"title"`
			Subject []string `json:"subject"`
			Author  []struct {
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
			Link []struct {
				URL                 string `json:"URL"`
				ContentType         string `json:"content-type"`
				ContentVersion      string `json:"content-version"`
				IntendedApplication string `json:"intended-application"`
			} `json:"link"`
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

	// Enforce has-full-text for much better PDF links
	searchURL := fmt.Sprintf("%s?q=%s&filter=has-full-text:true,language:pt&rows=%d", h.BaseURL, url.QueryEscape(searchTerm), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("CAPES harvester: request failed", zap.Error(err))
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
	resultsChan := make(chan material.Material, len(data.Message.Items))
	var wg sync.WaitGroup

	for _, item := range data.Message.Items {
		if len(item.Title) == 0 {
			continue
		}

		wg.Add(1)
		go func(it struct {
			Title   []string `json:"title"`
			Subject []string `json:"subject"`
			Author  []struct {
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
			Link []struct {
				URL                 string `json:"URL"`
				ContentType         string `json:"content-type"`
				ContentVersion      string `json:"content-version"`
				IntendedApplication string `json:"intended-application"`
			} `json:"link"`
		}) {
			defer wg.Done()

			var authors []string
			for _, a := range it.Author {
				authors = append(authors, fmt.Sprintf("%s %s", a.Given, a.Family))
			}

			year := 0
			if len(it.Created.DateParts) > 0 && len(it.Created.DateParts[0]) > 0 {
				year = it.Created.DateParts[0][0]
			}

			catName := category
			if catName == "" {
				if len(it.Subject) > 0 {
					catName = it.Subject[0]
				} else {
					catName = "Artigo Periódico"
				}
			}

			var pdfURL string
			for _, link := range it.Link {
				lowerURL := strings.ToLower(link.URL)
				// Strict PDF check: must end in .pdf or have explicit application/pdf content type
				if strings.HasSuffix(lowerURL, ".pdf") || link.ContentType == "application/pdf" {
					pdfURL = link.URL
					break
				}
			}

			if pdfURL == "" {
				return
			}

			// Fetch cover from Google Books (In parallel now)
			cover := GetCoverFromGoogleBooks(it.Title[0], strings.Join(authors, ", "))

			difficulty := 3
			if len(it.Abstract) > 1000 {
				difficulty = 4
			}
			if year < 2010 {
				difficulty++
			}
			if difficulty > 5 {
				difficulty = 5
			}

			xp := 10 + (difficulty * 5)

			resultsChan <- material.Material{
				Titulo:        it.Title[0],
				Autor:         strings.Join(authors, ", "),
				Descricao:     it.Abstract,
				AnoPublicacao: year,
				Fonte:         "CAPES",
				Categoria:     catName,
				ExternoID:     it.DOI,
				CapaURL:       cover,
				PDFURL:        pdfURL,
				Disponivel:    true,
				Dificuldade:   difficulty,
				XP:            xp,
				Relevancia:    10,
			}
		}(item)
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for m := range resultsChan {
		materials = append(materials, m)
	}

	logger.Info("CAPES harvester: refined search completed", zap.Int("results", len(materials)))
	return materials, nil
}
