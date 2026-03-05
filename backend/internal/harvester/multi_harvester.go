package harvester

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/logger"
	"context"

	"go.uber.org/zap"
)

type MultiSourceHarvester struct {
	capes *CAPESHarvester
}

func NewMultiSourceHarvester() *MultiSourceHarvester {
	return &MultiSourceHarvester{
		capes: NewCAPESHarvester(),
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
			allMaterials = append(allMaterials, m)
		}
	}

	logger.Info("CAPES search completed", zap.Int("total_results", len(allMaterials)))
	return allMaterials, nil
}

