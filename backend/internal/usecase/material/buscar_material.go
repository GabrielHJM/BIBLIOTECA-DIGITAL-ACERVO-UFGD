package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/pkg/metadata"
	"context"
)

type BuscarMaterialUseCase struct {
	Repo material.Repository
	Meta *metadata.MetadataService
}

func (uc *BuscarMaterialUseCase) Execute(ctx context.Context, id int) (*material.Material, error) {
	m, err := uc.Repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Se faltar capa ou descrição, tenta enriquecer e salvar
	if m.CapaURL == "" || m.Descricao == "" {
		cover, desc := uc.Meta.FetchEnrichment(m.Titulo, m.Autor, m.ISBN)
		updated := false
		if m.CapaURL == "" && cover != "" {
			m.CapaURL = cover
			updated = true
		}
		if m.Descricao == "" && desc != "" {
			m.Descricao = desc
			updated = true
		}

		if updated {
			// Salva a atualização no banco silenciosamente
			go uc.Repo.Atualizar(context.Background(), m)
		}
	}

	return m, nil
}
