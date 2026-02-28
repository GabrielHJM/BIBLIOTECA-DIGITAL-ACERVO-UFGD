package social

import "time"

// Like represents a like on a material by a user.
type Like struct {
    ID         int       `json:"id"`
    UsuarioID  int       `json:"usuario_id"`
    MaterialID int       `json:"material_id"`
    CreatedAt  time.Time `json:"created_at"`
}

// Comentario represents a comment on a material.
type Comentario struct {
    ID         int       `json:"id"`
    UsuarioID  int       `json:"usuario_id"`
    MaterialID int       `json:"material_id"`
    Texto      string    `json:"texto"`
    CreatedAt  time.Time `json:"created_at"`
}

// ShareLog records an external share action.
type ShareLog struct {
    ID         int       `json:"id"`
    UsuarioID  int       `json:"usuario_id"`
    MaterialID int       `json:"material_id"`
    Plataforma string    `json:"plataforma"`
    SharedAt   time.Time `json:"shared_at"`
}

// Mensagem represents a private message sending a material to a friend.
type Mensagem struct {
    ID            int       `json:"id"`
    RemetenteID   int       `json:"remetente_id"`
    DestinatarioID int      `json:"destinatario_id"`
    MaterialID    int       `json:"material_id"`
    Texto         string    `json:"texto"`
    SentAt        time.Time `json:"sent_at"`
}

// UseCase defines the business logic for social interactions.
type UseCase interface {
	ToggleLike(usuarioID, materialID int) (bool, error)
	AddComment(c Comentario) error
	RegisterShare(s ShareLog) error
	SendMessage(m Mensagem) error
}
