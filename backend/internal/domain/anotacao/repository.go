package anotacao

import "context"

type Repository interface {
	Create(ctx context.Context, req Anotacao) (int, error)
	GetByID(ctx context.Context, id int) (Anotacao, error)
	ListByUsuario(ctx context.Context, usuarioID int) ([]Anotacao, error)
	Update(ctx context.Context, req Anotacao) error
	Delete(ctx context.Context, id int, usuarioID int) error
}
