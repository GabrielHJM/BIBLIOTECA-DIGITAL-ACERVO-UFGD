package harvester

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// GetCoverFromGoogleBooks fetches the highest quality cover available for a given title and author.
func GetCoverFromGoogleBooks(title string, author string) string {
	query := title
	if author != "" {
		query += " " + author
	}
	searchURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=1", url.QueryEscape(query))

	// Make a single rapid request
	resp, err := http.Get(searchURL)
	if err != nil {
		return abstractAcademicCover(title)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return abstractAcademicCover(title)
	}

	var data struct {
		Items []struct {
			VolumeInfo struct {
				ImageLinks struct {
					Thumbnail  string `json:"thumbnail"`
					Large      string `json:"large"`
					ExtraLarge string `json:"extraLarge"`
				} `json:"imageLinks"`
			} `json:"volumeInfo"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return abstractAcademicCover(title)
	}

	if len(data.Items) > 0 {
		img := data.Items[0].VolumeInfo.ImageLinks
		cover := img.ExtraLarge
		if cover == "" {
			cover = img.Large
		}
		if cover == "" {
			cover = img.Thumbnail
		}
		if cover != "" {
			return strings.ReplaceAll(cover, "http://", "https://")
		}
	}

	return abstractAcademicCover(title)
}

// abstractAcademicCover returns a dynamically generated cover with the book's title and author.
func abstractAcademicCover(title string) string {
	// Clean string for URL
	safeTitle := url.QueryEscape(title)
	// We use short version if it's too long
	if len(title) > 30 {
		safeTitle = url.QueryEscape(title[:27] + "...")
	}

	backgrounds := []string{"1e1e1e", "0f172a", "2d3748", "111827"}
	textColors := []string{"00B8D4", "38bdf8", "a78bfa", "f472b6"}

	// Pseudo-random selection based on title length
	idx := len(title) % len(backgrounds)
	bg := backgrounds[idx]
	color := textColors[idx]

	return fmt.Sprintf("https://placehold.co/400x600/%s/%s?text=%s", bg, color, safeTitle)
}
