package material

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
)

type FavoritarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *FavoritarMaterialUseCase) Execute(ctx context.Context, usuarioID, materialID int, favoritar bool) error {
	if favoritar {
		f := &material.Favorito{UsuarioID: usuarioID, MaterialID: materialID}
		return uc.Repo.AdicionarFavorito(ctx, f)
	}
	return uc.Repo.RemoverFavorito(ctx, usuarioID, materialID)
}

func (uc *FavoritarMaterialUseCase) Listar(ctx context.Context, usuarioID int) ([]material.Material, error) {
	return uc.Repo.ListarFavoritosPorUsuario(ctx, usuarioID)
}
