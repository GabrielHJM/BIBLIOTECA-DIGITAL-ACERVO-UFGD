package social

// LikeRepository defines methods for persisting likes.
type LikeRepository interface {
	AddLike(usuarioID, materialID int) error
	RemoveLike(usuarioID, materialID int) error
	CountLikes(materialID int) (int, error)
	HasLiked(usuarioID, materialID int) (bool, error)
}

// ComentarioRepository defines methods for comments.
type ComentarioRepository interface {
	AddComentario(c Comentario) error
	ListComentarios(materialID int) ([]Comentario, error)
}

// ShareLogRepository defines methods for share logs.
type ShareLogRepository interface {
	AddShareLog(s ShareLog) error
	ListShares(materialID int) ([]ShareLog, error)
}

// MensagemRepository defines methods for private messages.
type MensagemRepository interface {
	SendMensagem(m Mensagem) error
	ListMensagens(usuarioID int) ([]Mensagem, error)
}
