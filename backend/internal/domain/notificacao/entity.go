package notificacao

import (
	"time"
)

type Notificacao struct {
	ID          int       `json:"id"`
	UsuarioID   int       `json:"usuario_id"`
	Titulo      string    `json:"titulo"`
	Mensagem    string    `json:"mensagem"`
	Tipo        string    `json:"tipo"`
	Lida        bool      `json:"lida"`
	DataCriacao time.Time `json:"data_criacao"`
}
