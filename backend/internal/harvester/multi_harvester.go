package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/pkg/metadata"
	"context"

	"go.uber.org/zap"
)

type MultiSourceHarvester struct {
	capes *CAPESHarvester
	meta  *metadata.MetadataService
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		capes: NewCAPESHarvester(),
		meta:  metadata.NewMetadataService(),
	}
}

func (h *MultiSourceHarvester) Search(ctx context.Context, query string, category string, source string, startYear int, endYear int, limit int) ([]material.Material, error) {
	// Focusing only on CAPES as requested
	mats, err := h.capes.Search(ctx, query, category, limit)
	if err != nil {
		logger.Error("CAPES Harvester failed", zap.Error(err))
		return nil, err
	}

	var allMaterials []material.Material
	seen := make(map[string]bool)

	for _, m := range mats {
		id := m.ExternoID
		if id == "" {
			id = m.Titulo
		}
		if !seen[id] {
			seen[id] = true
			// Enrich metadata if missing cover or description
			if m.CapaURL == "" || m.Descricao == "" {
				cover, desc := h.meta.FetchEnrichment(m.Titulo, m.Autor, m.ISBN)
				if m.CapaURL == "" && cover != "" {
					m.CapaURL = cover
				}
				if m.Descricao == "" && desc != "" {
					m.Descricao = desc
				}
			}
			allMaterials = append(allMaterials, m)
		}
	}

	logger.Info("CAPES search completed", zap.Int("total_results", len(allMaterials)))
	return allMaterials, nil
}

