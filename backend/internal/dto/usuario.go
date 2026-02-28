package dto

type UsuarioRequest struct {
	Nome    string `json:"nome"`
	Email   string `json:"email"`
	Senha   string `json:"senha"`
	Tipo    int    `json:"tipo"`
	FotoURL string `json:"foto_url"`
}

type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type RedefinirSenhaRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

type AtualizarMetaRequest struct {
	MetaPaginasSemana int `json:"meta_paginas_semana"`
}
