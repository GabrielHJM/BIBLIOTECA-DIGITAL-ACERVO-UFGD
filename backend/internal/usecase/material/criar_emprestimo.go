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
	// 1. Verifica se o material existe e está disponível
	m, err := uc.Repo.BuscarPorID(ctx, materialID)
	if err != nil {
		return err
	}
	if !m.Disponivel {
		return fmt.Errorf("material não disponível para empréstimo")
	}

	// 2. Cria o registro de empréstimo
	emprestimo := &material.Emprestimo{
		UsuarioID:      usuarioID,
		MaterialID:     materialID,
		DataEmprestimo: time.Now(),
		DataDevolucao:  time.Now().AddDate(0, 0, 7), // 7 dias
		Status:         "ativo",
	}

	if err := uc.Repo.SalvarEmprestimo(ctx, emprestimo); err != nil {
		return err
	}

	// 3. Atualiza status do material
	m.Disponivel = false
	return uc.Repo.Atualizar(ctx, m)
}
