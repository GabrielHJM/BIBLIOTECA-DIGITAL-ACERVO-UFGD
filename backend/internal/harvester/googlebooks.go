package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

type GoogleBooksResponse struct {
	Items []struct {
		ID         string `json:"id"`
		VolumeInfo struct {
			Title         string   `json:"title"`
			Authors       []string `json:"authors"`
			PublishedDate string   `json:"publishedDate"`
			Description   string   `json:"description"`
			PageCount     int      `json:"pageCount"`
			Categories    []string `json:"categories"`
			ImageLinks    struct {
				Thumbnail string `json:"thumbnail"`
			} `json:"imageLinks"`
			PreviewLink string `json:"previewLink"`
			InfoLink    string `json:"infoLink"`
		} `json:"volumeInfo"`
		AccessInfo struct {
			Viewability string `json:"viewability"`
			Pdf         struct {
				IsAvailable  bool   `json:"isAvailable"`
				DownloadLink string `json:"downloadLink"`
				AcsTokenLink string `json:"acsTokenLink"`
			} `json:"pdf"`
			WebReaderLink string `json:"webReaderLink"`
		} `json:"accessInfo"`
	} `json:"items"`
}

type GoogleBooksHarvester struct {
	BaseURL string
}

func NewGoogleBooksHarvester() *GoogleBooksHarvester {
	return &GoogleBooksHarvester{
		BaseURL: "https://www.googleapis.com/books/v1/volumes",
	}
}

func (h *GoogleBooksHarvester) Search(ctx context.Context, query string, category string, limit int) ([]material.Material, error) {
	limiter := GetRateLimiter()

	searchTerm := query
	if searchTerm == "" {
		searchTerm = category
	}
	if searchTerm == "" {
		searchTerm = "science"
	}

	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	keyParam := ""
	if apiKey != "" {
		keyParam = "&key=" + apiKey
	}
	langParam := "&langRestrict=pt"
	
	searchURL := fmt.Sprintf("%s?q=%s&filter=free-ebooks%s&maxResults=%d%s", h.BaseURL, url.QueryEscape(searchTerm), langParam, limit, keyParam)

	// Max 3 retries with exponential backoff
	var resp *http.Response
	var err error
	backoff := 500 * time.Millisecond

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (Google Books: ~10 requests per second is safe)
		limiter.Wait(ctx, ProviderGoogleBooks, 10, 5)

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
				logger.Warn("GoogleBooks rate limit hit, retrying...", zap.Int("attempt", i+1))
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
			return nil, fmt.Errorf("googlebooks api error: %s", resp.Status)
		}

		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			logger.Error("GoogleBooks harvester: request failed", zap.Error(err))
			return nil, err
		}
		return nil, fmt.Errorf("googlebooks api error after retries")
	}
	defer resp.Body.Close()

	var data GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var materials []material.Material
	for _, item := range data.Items {
		if item.VolumeInfo.Title == "" || item.AccessInfo.Viewability == "NO_PAGES" {
			continue
		}

		pdfURL := ""
		// Prioritize WebReaderLink for "online reading" in the browser
		if item.AccessInfo.WebReaderLink != "" {
			pdfURL = item.AccessInfo.WebReaderLink
		} else if item.AccessInfo.Pdf.IsAvailable {
			if item.AccessInfo.Pdf.DownloadLink != "" {
				pdfURL = item.AccessInfo.Pdf.DownloadLink
			} else if item.AccessInfo.Pdf.AcsTokenLink != "" {
				pdfURL = item.AccessInfo.Pdf.AcsTokenLink
			}
		}

		if pdfURL == "" {
			if item.VolumeInfo.PreviewLink != "" {
				pdfURL = item.VolumeInfo.PreviewLink
			}
		}

		if pdfURL == "" {
			continue
		}

		year := 0
		if len(item.VolumeInfo.PublishedDate) >= 4 {
			fmt.Sscanf(item.VolumeInfo.PublishedDate[:4], "%d", &year)
		}

		cover := item.VolumeInfo.ImageLinks.Thumbnail
		if cover != "" {
			cover = strings.ReplaceAll(cover, "http://", "https://")
		} else {
			cover = abstractAcademicCover(item.VolumeInfo.Title)
		}

		difficulty := 2
		if item.VolumeInfo.PageCount > 300 {
			difficulty = 3
		}
		if item.VolumeInfo.PageCount > 600 {
			difficulty = 4
		}

		xp := 10 + (difficulty * 5)
		relevance := 30 // High relevance for books

		cat := category
		if cat == "" {
			if len(item.VolumeInfo.Categories) > 0 {
				cat = item.VolumeInfo.Categories[0]
			} else {
				cat = "Livro"
			}
		}

		m := material.Material{
			Titulo:        item.VolumeInfo.Title,
			Autor:         strings.Join(item.VolumeInfo.Authors, ", "),
			Descricao:     item.VolumeInfo.Description,
			AnoPublicacao: year,
			Paginas:       item.VolumeInfo.PageCount,
			Fonte:         "Google Books",
			Categoria:     cat,
			ExternoID:     item.ID,
			CapaURL:       cover,
			PDFURL:        pdfURL,
			Disponivel:    true,
			Dificuldade:   difficulty,
			XP:            xp,
			Relevancia:    relevance,
		}

		materials = append(materials, m)
	}

	logger.Info("GoogleBooks harvester: search completed", zap.Int("results", len(materials)))
	return materials, nil
}
