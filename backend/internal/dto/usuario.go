package dto

type UsuarioRequest struct {
	Nome  string `json:"nome"`
	Email string `json:"email"`
	Senha string `json:"senha"`
	Tipo  int    `json:"tipo"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type RedefinirSenhaRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}
