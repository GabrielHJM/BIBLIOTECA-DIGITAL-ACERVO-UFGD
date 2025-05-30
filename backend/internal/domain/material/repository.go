package material

type Repository interface {
    Listar() ([]Material, error)
    Favoritar(id int) error
}
