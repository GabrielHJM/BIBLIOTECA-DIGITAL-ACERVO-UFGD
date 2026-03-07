package anotacao

import "time"

type Anotacao struct {
	ID              int       `json:"id"`
	UsuarioID       int       `json:"usuario_id"`
	MaterialID      *int      `json:"material_id,omitempty"`
	Titulo          string    `json:"titulo"`
	Conteudo        string    `json:"conteudo"`
	Cor             string    `json:"cor"`
	DataCriacao     time.Time `json:"data_criacao"`
	DataAtualizacao time.Time `json:"data_atualizacao"`
}
