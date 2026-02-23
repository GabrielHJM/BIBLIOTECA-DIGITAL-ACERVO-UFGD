package material

import "time"

type Emprestimo struct {
	ID             int       `json:"id"`
	UsuarioID      int       `json:"usuario_id"`
	MaterialID     int       `json:"material_id"`
	DataEmprestimo time.Time `json:"data_emprestimo"`
	DataDevolucao  time.Time `json:"data_devolucao"`
	Status         string    `json:"status"` // 'pendente', 'devolvido', 'atrasado'
}

type Avaliacao struct {
	ID         int       `json:"id"`
	UsuarioID  int       `json:"usuario_id"`
	MaterialID int       `json:"material_id"`
	Nota       int       `json:"nota"` // 1-5
	Comentario string    `json:"comentario"`
	Data       time.Time `json:"data"`
}

type Favorito struct {
	UsuarioID  int `json:"usuario_id"`
	MaterialID int `json:"material_id"`
}

type HistoricoLeitura struct {
	ID         int       `json:"id"`
	UsuarioID  int       `json:"usuario_id"`
	MaterialID int       `json:"material_id"`
	Data       time.Time `json:"data"`
}
