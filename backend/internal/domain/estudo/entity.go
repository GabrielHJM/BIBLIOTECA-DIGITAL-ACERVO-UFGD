package estudo

import (
	"time"
)

type Anotacao struct {
	ID          int       `json:"id"`
	UsuarioID   int       `json:"usuario_id"`
	MaterialID  int       `json:"material_id"`
	Conteudo    string    `json:"conteudo"`
	Pagina      int       `json:"pagina"`
	Cor         string    `json:"cor"`
	DataCriacao time.Time `json:"data_criacao"`
}

type Flashcard struct {
	ID             int       `json:"id"`
	UsuarioID      int       `json:"usuario_id"`
	MaterialID     int       `json:"material_id"`
	Pergunta       string    `json:"pergunta"`
	Resposta       string    `json:"resposta"`
	Dificuldade    int       `json:"dificuldade"`
	ProximaRevisao time.Time `json:"proxima_revisao"`
	DataCriacao    time.Time `json:"data_criacao"`
}
