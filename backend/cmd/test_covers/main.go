package main

import (
	"biblioteca-digital-api/internal/harvester"
	"fmt"
)

func main() {
	testCases := []struct {
		Title  string
		Author string
	}{
		{"Clean Code", "Robert C. Martin"},
		{"Ciencia-Tecnología-Sociedad vs. STEM", ""},
		{"Orchestrating Intelligence", "Usman Durrani"},
	}

	for _, tc := range testCases {
		fmt.Printf("Buscando capa para: '%s' por '%s'\n", tc.Title, tc.Author)
		cover := harvester.GetCoverFromGoogleBooks(tc.Title, tc.Author)
		if cover != "" {
			fmt.Printf("-> Capa retornada: %s\n", cover)
		} else {
			fmt.Println("-> FALHA GERAL")
		}
		fmt.Println("---")
	}
}
