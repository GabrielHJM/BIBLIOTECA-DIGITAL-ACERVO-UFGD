package material

import "biblioteca-digital-api/internal/domain/material"

type FavoritarMaterialUseCase struct {
	Repo material.Repository
}

func (uc *FavoritarMaterialUseCase) Execute(usuarioID, materialID int, favoritar bool) error {
	if favoritar {
		f := &material.Favorito{
			UsuarioID:  usuarioID,
			MaterialID: materialID,
		}
		return uc.Repo.AdicionarFavorito(f)
	}
	return uc.Repo.RemoverFavorito(usuarioID, materialID)
}

func (uc *FavoritarMaterialUseCase) Listar(usuarioID int) ([]material.Material, error) {
	return uc.Repo.ListarFavoritosPorUsuario(usuarioID)
}
