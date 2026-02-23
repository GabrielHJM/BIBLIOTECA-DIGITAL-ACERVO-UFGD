package usuario

type Usuario struct {
	ID         int
	Nome       string
	Email      string
	Senha      string
	Tipo       int
	Interesses []string
}
