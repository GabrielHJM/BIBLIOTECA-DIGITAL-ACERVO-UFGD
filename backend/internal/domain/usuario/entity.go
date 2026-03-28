package usuario

import "time"

type Usuario struct {
	ID             int
	Nome           string
	Email          string
	Senha          string
	Tipo           int
	Interesses     []string
	FotoURL        string
	Cpf            string
	DataNascimento *time.Time
	Username       string
	DeletedAt      *time.Time
}
