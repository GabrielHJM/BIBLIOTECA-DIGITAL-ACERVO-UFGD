package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type GoogleBooksResponse struct {
	Items []struct {
		VolumeInfo struct {
			Title       string   `json:"title"`
			Authors     []string `json:"authors"`
			Description string   `json:"description"`
			ImageLinks  struct {
				Thumbnail      string `json:"thumbnail"`
				SmallThumbnail string `json:"smallThumbnail"`
			} `json:"imageLinks"`
		} `json:"volumeInfo"`
	} `json:"items"`
}

type MetadataService struct {
	client *http.Client
}

func NewMetadataService() *MetadataService {
	return &MetadataService{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

// FetchEnrichment tenta buscar capa e resumo baseado no título e autor ou ISBN
func (s *MetadataService) FetchEnrichment(title, author, isbn string) (coverURL, description string) {
	if title == "" && isbn == "" {
		return "", ""
	}

	query := ""
	if isbn != "" {
		query = "isbn:" + isbn
	} else {
		query = title
		if author != "" {
			query += " " + author
		}
	}

	apiURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=1", url.QueryEscape(query))

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", ""
	}

	var data GoogleBooksResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", ""
	}

	if len(data.Items) > 0 {
		info := data.Items[0].VolumeInfo
		cover := info.ImageLinks.Thumbnail
		// Upgrade thumbnail to https and larger version if possible
		if strings.HasPrefix(cover, "http://") {
			cover = strings.Replace(cover, "http://", "https://", 1)
		}
		// Algumas vezes a Google API retorna capas pequenas, tentamos forçar zoom se disponível
		// (Isso é um hack comum para Google Books)

		return cover, info.Description
	}

	return "", ""
}
