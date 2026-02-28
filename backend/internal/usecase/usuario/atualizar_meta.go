package usuario

import (
	"biblioteca-digital-api/internal/domain/usuario"
	"context"
)

type AtualizarMetaUseCase struct {
	repo usuario.UsuarioRepository
}

func NewAtualizarMeta(repo usuario.UsuarioRepository) *AtualizarMetaUseCase {
	return &AtualizarMetaUseCase{repo}
}

func (uc *AtualizarMetaUseCase) Execute(ctx context.Context, id int, meta int) error {
	return uc.repo.AtualizarMeta(ctx, id, meta)
}
