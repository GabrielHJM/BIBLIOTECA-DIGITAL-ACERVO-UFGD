package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"fmt"
	"time"
)

type CriarEmprestimoUseCase struct {
	Repo material.Repository
}

func (uc *CriarEmprestimoUseCase) Execute(ctx context.Context, usuarioID, materialID int) error {
	m, err := uc.Repo.BuscarPorID(materialID)
	if err != nil {
		return err
	}
	if !m.Disponivel {
		return fmt.Errorf("material não disponível para empréstimo")
	}

	emp := &material.Emprestimo{
		UsuarioID:      usuarioID,
		MaterialID:     materialID,
		DataEmprestimo: time.Now(),
		DataDevolucao:  time.Now().AddDate(0, 0, 14), // 14 dias
		Status:         "pendente",
	}

	if err := uc.Repo.SalvarEmprestimo(emp); err != nil {
		return err
	}

	// Atualiza disponibilidade
	m.Disponivel = false
	return uc.Repo.Atualizar(m)
}
