package utils

import (
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// RemoveAccents removes diacritics from a string.
func RemoveAccents(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

// SortByFieldReverse is a helper to sort any slice which has a Score field.
// Since we have a specific case in the usecase, we can just implement it there or here with a generic interface.
// For simplicity and speed, let's provide a generic-like way if possible, or just the specific one.
type Scored interface {
	GetScore() float64
}

func SortScoredReverse[T Scored](slice []T) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].GetScore() > slice[j].GetScore()
	})
}

// LevenshteinDistance calculates the Levenshtein distance between two strings.
func LevenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	if n == 0 { return m }
	if m == 0 { return n }
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
		d[i][0] = i
	}
	for j := 0; j <= m; j++ { d[0][j] = j }
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] { cost = 0 }
			d[i][j] = min(d[i-1][j]+1, min(d[i][j-1]+1, d[i-1][j-1]+cost))
		}
	}
	return d[n][m]
}

// IsTitleSimilar returns true if titles are >= 85% similar.
func IsTitleSimilar(t1, t2 string) bool {
	t1 = strings.ToLower(RemoveAccents(strings.TrimSpace(t1)))
	t2 = strings.ToLower(RemoveAccents(strings.TrimSpace(t2)))
	if t1 == t2 { return true }
	if len(t1) == 0 || len(t2) == 0 { return false }

	dist := LevenshteinDistance(t1, t2)
	maxLen := max(len(t1), len(t2))
	similarity := 1.0 - float64(dist)/float64(maxLen)
	
	// Consider similar if > 85% match or if one string is fully contained in another and lengths are close
	if similarity >= 0.85 { return true }
	if (strings.Contains(t1, t2) || strings.Contains(t2, t1)) && similarity >= 0.60 { return true }
	
	return false
}

// ParseSmartQuery extracts year from a query string and returns the cleaned query and the year.
func ParseSmartQuery(query string) (string, int) {
	if len(strings.TrimSpace(query)) < 4 {
		return query, 0
	}
	re := regexp.MustCompile(`\b(19\d{2}|20[0-2]\d)\b`)
	match := re.FindString(query)
	if match != "" {
		year, _ := strconv.Atoi(match)
		cleanQuery := strings.TrimSpace(re.ReplaceAllString(query, ""))
		// Limpa multiplos espaços gerados
		cleanQuery = regexp.MustCompile(`\s+`).ReplaceAllString(cleanQuery, " ")
		return cleanQuery, year
	}
	return query, 0
}
