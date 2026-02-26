package main

import (
	"biblioteca-digital-api/config"
	"biblioteca-digital-api/internal/domain/material"
	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"
	"context"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		logger.Info("Nenhum arquivo .env encontrado")
	}

	cfg := config.Load()
	db := config.InitDB(cfg)
	repo := &repository.MaterialPostgres{DB: db}

	source := "all"
	if len(os.Args) > 1 {
		source = os.Args[1]
	}

	logger.Info("Iniciando ingestão de dados", zap.String("source", source))

	if source == "all" || source == "scielo" {
		h := harvester.NewSciELOHarvester()
		materials, err := h.Search(context.Background(), "", "", 20)
		if err != nil {
			logger.Error("Erro ao coletar SciELO", zap.Error(err))
		} else {
			saveMaterials(repo, materials)
		}
	}

	if source == "all" || source == "ieee" {
		h := harvester.NewIEEEHarvester()
		materials, err := h.Search(context.Background(), "tecnologia", "", 10)
		if err != nil {
			logger.Error("Erro ao coletar IEEE", zap.Error(err))
		} else {
			saveMaterials(repo, materials)
		}
	}

	if source == "all" || source == "capes" {
		h := harvester.NewCAPESHarvester()
		materials, err := h.Search(context.Background(), "ciência", "", 10)
		if err != nil {
			logger.Error("Erro ao coletar CAPES", zap.Error(err))
		} else {
			saveMaterials(repo, materials)
		}
	}

	logger.Info("Ingestão concluída")
}

func saveMaterials(repo *repository.MaterialPostgres, materials []material.Material) {
	for i := range materials {
		m := &materials[i]
		if m.Categoria == "" {
			m.Categoria = "Acadêmico"
		}

		err := repo.Criar(context.Background(), m)
		if err != nil {
			logger.Error("Erro ao salvar material", zap.String("titulo", m.Titulo), zap.Error(err))
		} else {
			logger.Info("Material salvo", zap.String("titulo", m.Titulo), zap.Int("id", m.ID))
		}
	}
}
