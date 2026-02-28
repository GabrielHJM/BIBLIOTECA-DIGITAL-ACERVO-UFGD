package repository

import (
	"biblioteca-digital-api/internal/domain/social"
	"database/sql"
	"time"
)

type SocialPostgres struct {
	DB *sql.DB
}

func NewSocialPG(db *sql.DB) *SocialPostgres {
	return &SocialPostgres{DB: db}
}

// LikeRepository implementation
func (r *SocialPostgres) AddLike(usuarioID, materialID int) error {
	query := `INSERT INTO likes (usuario_id, material_id, created_at) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.DB.Exec(query, usuarioID, materialID, time.Now())
	return err
}

func (r *SocialPostgres) RemoveLike(usuarioID, materialID int) error {
	query := `DELETE FROM likes WHERE usuario_id = $1 AND material_id = $2`
	_, err := r.DB.Exec(query, usuarioID, materialID)
	return err
}

func (r *SocialPostgres) CountLikes(materialID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM likes WHERE material_id = $1`
	err := r.DB.QueryRow(query, materialID).Scan(&count)
	return count, err
}

func (r *SocialPostgres) HasLiked(usuarioID, materialID int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE usuario_id = $1 AND material_id = $2)`
	err := r.DB.QueryRow(query, usuarioID, materialID).Scan(&exists)
	return exists, err
}

// ComentarioRepository implementation
func (r *SocialPostgres) AddComentario(c social.Comentario) error {
	query := `INSERT INTO comentarios (usuario_id, material_id, texto, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, c.UsuarioID, c.MaterialID, c.Texto, time.Now())
	return err
}

func (r *SocialPostgres) ListComentarios(materialID int) ([]social.Comentario, error) {
	rows, err := r.DB.Query(`SELECT id, usuario_id, material_id, texto, created_at FROM comentarios WHERE material_id = $1 ORDER BY created_at DESC`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comentarios []social.Comentario
	for rows.Next() {
		var c social.Comentario
		if err := rows.Scan(&c.ID, &c.UsuarioID, &c.MaterialID, &c.Texto, &c.CreatedAt); err != nil {
			return nil, err
		}
		comentarios = append(comentarios, c)
	}
	return comentarios, nil
}

// ShareLogRepository implementation
func (r *SocialPostgres) AddShareLog(s social.ShareLog) error {
	query := `INSERT INTO share_logs (usuario_id, material_id, plataforma, shared_at) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, s.UsuarioID, s.MaterialID, s.Plataforma, time.Now())
	return err
}

func (r *SocialPostgres) ListShares(materialID int) ([]social.ShareLog, error) {
	rows, err := r.DB.Query(`SELECT id, usuario_id, material_id, plataforma, shared_at FROM share_logs WHERE material_id = $1 ORDER BY shared_at DESC`, materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var shares []social.ShareLog
	for rows.Next() {
		var s social.ShareLog
		if err := rows.Scan(&s.ID, &s.UsuarioID, &s.MaterialID, &s.Plataforma, &s.SharedAt); err != nil {
			return nil, err
		}
		shares = append(shares, s)
	}
	return shares, nil
}

// MensagemRepository implementation
func (r *SocialPostgres) SendMensagem(m social.Mensagem) error {
	query := `INSERT INTO mensagens (remetente_id, destinatario_id, material_id, texto, sent_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, m.RemetenteID, m.DestinatarioID, m.MaterialID, m.Texto, time.Now())
	return err
}

func (r *SocialPostgres) ListMensagens(usuarioID int) ([]social.Mensagem, error) {
	rows, err := r.DB.Query(`SELECT id, remetente_id, destinatario_id, material_id, texto, sent_at FROM mensagens WHERE destinatario_id = $1 ORDER BY sent_at DESC`, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var msgs []social.Mensagem
	for rows.Next() {
		var m social.Mensagem
		if err := rows.Scan(&m.ID, &m.RemetenteID, &m.DestinatarioID, &m.MaterialID, &m.Texto, &m.SentAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, nil
}
