package harvester

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetCoverFromGoogleBooks fetches the highest quality cover available for a given title and author.
func GetCoverFromGoogleBooks(title string, author string) string {
	ctx := context.Background()
	limiter := GetRateLimiter()

	query := title
	if author != "" {
		query += " " + author
	}
	searchURL := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s&maxResults=1", url.QueryEscape(query))

	// Max 3 retries with exponential backoff
	var resp *http.Response
	var err error
	backoff := 500 * time.Millisecond

	for i := 0; i < 3; i++ {
		// Wait for rate limiter (Google Books: ~10 requests per second to be safe)
		limiter.Wait(ctx, ProviderGoogleBooks, 10, 5)

		resp, err = http.Get(searchURL)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				break
			}
			if resp.StatusCode == http.StatusTooManyRequests {
				resp.Body.Close()
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			resp.Body.Close()
		}
		
		if i < 2 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	if err != nil || resp == nil || resp.StatusCode != http.StatusOK {
		return abstractAcademicCover(title)
	}
	defer resp.Body.Close()

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
