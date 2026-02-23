package repository

import (
	"biblioteca-digital-api/internal/domain/material"
	"database/sql"
)

type MaterialPostgres struct {
	DB *sql.DB
}

func (r *MaterialPostgres) Listar(limit, offset int) ([]material.Material, error) {
	query := "SELECT id, titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, disponivel, media_nota, total_avaliacoes FROM materiais LIMIT $1 OFFSET $2"
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) BuscarPorID(id int) (*material.Material, error) {
	var m material.Material
	err := r.DB.QueryRow("SELECT id, titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, disponivel, media_nota, total_avaliacoes FROM materiais WHERE id = $1", id).
		Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MaterialPostgres) Pesquisar(termo string, categoria string, tags []string, limit, offset int) ([]material.Material, error) {
	query := `SELECT id, titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, disponivel, media_nota, total_avaliacoes
			  FROM materiais
			  WHERE (titulo ILIKE $1 OR autor ILIKE $2 OR descricao ILIKE $3)`

	args := []interface{}{"%" + termo + "%", "%" + termo + "%", "%" + termo + "%"}
	argCount := 4

	if categoria != "" {
		query += " AND categoria = $" + string(rune(argCount+'0'))
		args = append(args, categoria)
		argCount++
	}

	// Simplificado: tags logic would require a separate table or jsonb
	// For now, let's keep it simple with limit/offset
	query += " LIMIT $" + string(rune(argCount+'0')) + " OFFSET $" + string(rune(argCount+1+'0'))
	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) Criar(m *material.Material) error {
	return r.DB.QueryRow(
		"INSERT INTO materiais (titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, disponivel) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id",
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.Disponivel,
	).Scan(&m.ID)
}

func (r *MaterialPostgres) Atualizar(m *material.Material) error {
	_, err := r.DB.Exec(
		"UPDATE materiais SET titulo=$1, autor=$2, isbn=$3, categoria=$4, ano_publicacao=$5, descricao=$6, capa_url=$7, disponivel=$8 WHERE id=$9",
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.Disponivel, m.ID,
	)
	return err
}

func (r *MaterialPostgres) Deletar(id int) error {
	_, err := r.DB.Exec("DELETE FROM materiais WHERE id=$1", id)
	return err
}

func (r *MaterialPostgres) SalvarEmprestimo(e *material.Emprestimo) error {
	return r.DB.QueryRow(
		"INSERT INTO emprestimos (usuario_id, material_id, data_emprestimo, data_devolucao, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		e.UsuarioID, e.MaterialID, e.DataEmprestimo, e.DataDevolucao, e.Status,
	).Scan(&e.ID)
}

func (r *MaterialPostgres) ListarEmprestimosPorUsuario(usuarioID int) ([]material.Emprestimo, error) {
	rows, err := r.DB.Query("SELECT id, usuario_id, material_id, data_emprestimo, data_devolucao, status FROM emprestimos WHERE usuario_id = $1", usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emprestimos []material.Emprestimo
	for rows.Next() {
		var e material.Emprestimo
		if err := rows.Scan(&e.ID, &e.UsuarioID, &e.MaterialID, &e.DataEmprestimo, &e.DataDevolucao, &e.Status); err != nil {
			return nil, err
		}
		emprestimos = append(emprestimos, e)
	}
	return emprestimos, nil
}
func (r *MaterialPostgres) SalvarAvaliacao(a *material.Avaliacao) error {
	return r.DB.QueryRow(
		"INSERT INTO avaliacoes (usuario_id, material_id, nota, comentario, data) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		a.UsuarioID, a.MaterialID, a.Nota, a.Comentario, a.Data,
	).Scan(&a.ID)
}

func (r *MaterialPostgres) ListarAvaliacoesPorMaterial(materialID int) ([]material.Avaliacao, error) {
	rows, err := r.DB.Query("SELECT id, usuario_id, material_id, nota, comentario, data FROM avaliacoes WHERE material_id = $1", materialID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var avaliacoes []material.Avaliacao
	for rows.Next() {
		var a material.Avaliacao
		if err := rows.Scan(&a.ID, &a.UsuarioID, &a.MaterialID, &a.Nota, &a.Comentario, &a.Data); err != nil {
			return nil, err
		}
		avaliacoes = append(avaliacoes, a)
	}
	return avaliacoes, nil
}

func (r *MaterialPostgres) AdicionarFavorito(f *material.Favorito) error {
	_, err := r.DB.Exec("INSERT INTO favoritos (usuario_id, material_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", f.UsuarioID, f.MaterialID)
	return err
}

func (r *MaterialPostgres) RemoverFavorito(usuarioID, materialID int) error {
	_, err := r.DB.Exec("DELETE FROM favoritos WHERE usuario_id = $1 AND material_id = $2", usuarioID, materialID)
	return err
}

func (r *MaterialPostgres) ListarFavoritosPorUsuario(usuarioID int) ([]material.Material, error) {
	query := `
		SELECT m.id, m.titulo, m.autor, m.isbn, m.categoria, m.ano_publicacao, m.descricao, m.capa_url, m.disponivel
		FROM materiais m
		JOIN favoritos f ON m.id = f.material_id
		WHERE f.usuario_id = $1`
	rows, err := r.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) RegistrarLeitura(h *material.HistoricoLeitura) error {
	return r.DB.QueryRow(
		"INSERT INTO historico_leitura (usuario_id, material_id, data) VALUES ($1, $2, $3) RETURNING id",
		h.UsuarioID, h.MaterialID, h.Data,
	).Scan(&h.ID)
}

func (r *MaterialPostgres) ListarHistoricoPorUsuario(usuarioID int) ([]material.Material, error) {
	query := `
		SELECT DISTINCT m.id, m.titulo, m.autor, m.isbn, m.categoria, m.ano_publicacao, m.descricao, m.capa_url, m.disponivel, m.media_nota, m.total_avaliacoes
		FROM materiais m
		JOIN historico_leitura h ON m.id = h.material_id
		WHERE h.usuario_id = $1
		ORDER BY m.id DESC`
	rows, err := r.DB.Query(query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) ObterRecomendacoes(usuarioID int, limit int) ([]material.Material, error) {
	// Simple recommendation: match user interests (categories) or just top rated if no interests
	query := `
		SELECT DISTINCT m.id, m.titulo, m.autor, m.isbn, m.categoria, m.ano_publicacao, m.descricao, m.capa_url, m.disponivel, m.media_nota, m.total_avaliacoes
		FROM materiais m
		WHERE m.categoria IN (SELECT interesse FROM interesses_usuario WHERE usuario_id = $1)
		OR m.media_nota >= 4.0
		ORDER BY m.media_nota DESC, m.total_avaliacoes DESC
		LIMIT $2`
	rows, err := r.DB.Query(query, usuarioID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}
