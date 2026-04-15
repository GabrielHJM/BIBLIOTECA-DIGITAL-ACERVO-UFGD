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

type DOAJResponse struct {
	Results []struct {
		BibJson struct {
			Title   string `json:"title"`
			Year    string `json:"year"`
			Journal struct {
				Title string `json:"title"`
			} `json:"journal"`
			Author []struct {
				Name string `json:"name"`
			} `json:"author"`
			Link []struct {
				Url  string `json:"url"`
				Type string `json:"type"`
			} `json:"link"`
			Identifier []struct {
				Id   string `json:"id"`
				Type string `json:"type"`
			} `json:"identifier"`
		} `json:"bibjson"`
	} `json:"results"`
}

type DOAJHarvester struct {
	BaseURL string
}

func NewDOAJHarvester() *DOAJHarvester {
	return &DOAJHarvester{
		BaseURL: "https://doaj.org/api/search/articles/",
	}
}

func (h *DOAJHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	// We are filtering strictly for Portuguese because DOAJ is huge and has a lot of content in many languages.
	// We can use Lucene query syntax: Portuguese AND (query)
	doajQuery := fmt.Sprintf(`language:"Portuguese" AND (%s)`, searchTerm)

	searchURL := fmt.Sprintf("%s%s?pageSize=%d", h.BaseURL, url.PathEscape(doajQuery), limit)

	var resp *http.Response
	var err error
	backoff := 500 * time.Millisecond

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (DOAJ allows ~40 requests per minute without API key -> ~1.5 sec per request)
		limiter.Wait(ctx, "doaj", 1, 1)

		req, errReq := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
		if errReq != nil {
			return nil, errReq
		}

		resp, err = http.DefaultClient.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				resp.Body.Close()
				logger.Warn("DOAJ rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("doaj api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("DOAJ harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("doaj api error after retries")
	}
	defer resp.Body.Close()

	var data DOAJResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Results {
		if item.BibJson.Title == "" {
			continue
		}

		pdfURL := ""
		for _, link := range item.BibJson.Link {
			if link.Type == "fulltext" {
				pdfURL = link.Url
				break
			}
		}

		if pdfURL == "" {
			continue
		}

		var authors []string
		for _, a := range item.BibJson.Author {
			authors = append(authors, a.Name)
		}

		year := 0
		if item.BibJson.Year != "" {
			fmt.Sscanf(item.BibJson.Year, "%d", &year)
		}

		cover := abstractAcademicCover(item.BibJson.Title)
		difficulty := 3 // Academic articles are generally high difficulty
		xp := 25
		relevance := 25

		cat := category
		if cat == "" {
			cat = "Artigo Acadêmico"
		}

		m := material.Material{
			Titulo:        item.BibJson.Title,
			Autor:         strings.Join(authors, ", "),
			Descricao:     "Publicado na revista: " + item.BibJson.Journal.Title,
			AnoPublicacao: year,
			Paginas:       0, // usually unknown without PDF parsing
			Fonte:         "DOAJ",
			Categoria:     cat,
			ExternoID:     "",
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   difficulty,
			XP:            xp,
			Relevancia:    relevance,
		}

		for _, id := range item.BibJson.Identifier {
			if id.Type == "doi" {
				m.ExternoID = id.Id
				m.ISBN = id.Id // we can use ISBN field to store DOI since the DB supports it
				break
			}
		}
		
		if m.ExternoID == "" {
			m.ExternoID = "doaj_" + fmt.Sprintf("%d", time.Now().UnixNano())
		}

		materials = append(materials, m)
	}

	logger.Info("DOAJ harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
