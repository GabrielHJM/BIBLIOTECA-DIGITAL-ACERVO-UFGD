package anotacao

import (
	"context"
	"errors"
	"time"

	domain "biblioteca-digital-api/internal/domain/anotacao"
)

type Usecase struct {
	repo domain.Repository
}

func NewUsecase(repo domain.Repository) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Criar(ctx context.Context, req domain.Anotacao) (int, error) {
	if req.Conteudo == "" {
		return 0, errors.New("o conteúdo da anotação não pode ser vazio")
	}
	if req.UsuarioID == 0 {
		return 0, errors.New("usuário inválido")
	}
	if req.Cor == "" {
		req.Cor = "#FFFFFF"
	}
	return u.repo.Create(ctx, req)
}

func (u *Usecase) ListarPorUsuario(ctx context.Context, usuarioID int) ([]domain.Anotacao, error) {
	return u.repo.ListByUsuario(ctx, usuarioID)
}

func (u *Usecase) ObterPorID(ctx context.Context, id int) (domain.Anotacao, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *Usecase) Atualizar(ctx context.Context, req domain.Anotacao) error {
	if req.ID == 0 || req.UsuarioID == 0 {
		return errors.New("id ou usuário inválido")
	}

	// Validar existência
	existente, err := u.repo.GetByID(ctx, req.ID)
	if err != nil {
		return errors.New("anotação não encontrada")
	}
	if existente.UsuarioID != req.UsuarioID {
		return errors.New("permissão negada")
	}

	req.DataAtualizacao = time.Now()
	return u.repo.Update(ctx, req)
}

func (u *Usecase) Excluir(ctx context.Context, id int, usuarioID int) error {
	return u.repo.Delete(ctx, id, usuarioID)
}
