package dto

type UsuarioRequest struct {
	Nome           string `json:"nome"`
	Email          string `json:"email"`
	Senha          string `json:"senha"` // For creation
	Tipo           int    `json:"tipo"`  // For creation
	FotoURL        string `json:"foto_url"`
	Cpf            string `json:"cpf"`
	DataNascimento string `json:"data_nascimento"`
	Username       string `json:"username"`
	SenhaAtual     string `json:"senha_atual"`
	NovaSenha      string `json:"nova_senha"`
}

type LoginRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required"`
}

type RedefinirSenhaRequest struct {
	Email string `json:"email" validate:"required,email"`
	Senha string `json:"senha" validate:"required,min=6"`
}

type AtualizarMetaRequest struct {
	MetaPaginasSemana int `json:"meta_paginas_semana"`
}
