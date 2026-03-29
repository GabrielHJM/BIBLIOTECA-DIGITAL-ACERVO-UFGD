package utils

import (
	"sort"
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
