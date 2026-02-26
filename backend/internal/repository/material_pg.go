package repository

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"database/sql"
	"fmt"
)

type MaterialPostgres struct {
	DB *sql.DB
}

func (r *MaterialPostgres) Listar(ctx context.Context, limit, offset int) ([]material.Material, error) {
	query := `SELECT id, titulo, autor, COALESCE(isbn, ''), categoria, ano_publicacao,
	          COALESCE(descricao, ''), COALESCE(capa_url, ''), COALESCE(pdf_url, ''),
	          disponivel, COALESCE(media_nota, 0.0), COALESCE(total_avaliacoes, 0)
	          FROM materiais LIMIT $1 OFFSET $2`
	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) BuscarPorID(ctx context.Context, id int) (*material.Material, error) {
	var m material.Material
	query := `SELECT id, titulo, autor, COALESCE(isbn, ''), categoria, ano_publicacao,
	          COALESCE(descricao, ''), COALESCE(capa_url, ''), COALESCE(pdf_url, ''),
	          disponivel, COALESCE(media_nota, 0.0), COALESCE(total_avaliacoes, 0)
	          FROM materiais WHERE id = $1`
	err := r.DB.QueryRowContext(ctx, query, id).
		Scan(&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MaterialPostgres) Pesquisar(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
	query := `SELECT id, titulo, autor, COALESCE(isbn, ''), categoria, ano_publicacao,
	          COALESCE(descricao, ''), COALESCE(capa_url, ''), COALESCE(pdf_url, ''),
	          disponivel, COALESCE(media_nota, 0.0), COALESCE(total_avaliacoes, 0)
			  FROM materiais
			  WHERE 1=1`

	args := []interface{}{}
	argCount := 1

	if termo != "" {
		query += fmt.Sprintf(" AND (titulo ILIKE $%d OR autor ILIKE $%d OR COALESCE(descricao, '') ILIKE $%d)", argCount, argCount, argCount)
		args = append(args, "%"+termo+"%")
		argCount++
	}

	if categoria != "" {
		query += fmt.Sprintf(" AND categoria ILIKE $%d", argCount)
		args = append(args, categoria)
		argCount++
	}

	if fonte != "" {
		query += fmt.Sprintf(" AND fonte ILIKE $%d", argCount)
		args = append(args, fonte)
		argCount++
	}

	if anoInicio > 0 {
		query += fmt.Sprintf(" AND ano_publicacao >= $%d", argCount)
		args = append(args, anoInicio)
		argCount++
	}

	if anoFim > 0 {
		query += fmt.Sprintf(" AND ano_publicacao <= $%d", argCount)
		args = append(args, anoFim)
		argCount++
	}

	orderBy := "id DESC"
	if sort == "random" {
		orderBy = "RANDOM()"
	}

	query += " ORDER BY " + orderBy
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) BuscarSimilares(ctx context.Context, materialID int, limit int) ([]material.Material, error) {
	query := `SELECT m2.id, m2.titulo, m2.autor, COALESCE(m2.isbn, ''), m2.categoria, m2.ano_publicacao,
	          COALESCE(m2.descricao, ''), COALESCE(m2.capa_url, ''), COALESCE(m2.pdf_url, ''),
	          m2.disponivel, COALESCE(m2.media_nota, 0.0), COALESCE(m2.total_avaliacoes, 0)
	          FROM materiais m1
	          JOIN materiais m2 ON m1.categoria = m2.categoria AND m1.id != m2.id
	          WHERE m1.id = $1
	          LIMIT $2`
	rows, err := r.DB.QueryContext(ctx, query, materialID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) Criar(ctx context.Context, m *material.Material) error {
	// O sistema verifica duplicatas no app porque ON CONFLICT externo_id falha sem o index
	if m.ExternoID != "" {
		var exists bool
		err := r.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM materiais WHERE externo_id = $1)", m.ExternoID).Scan(&exists)
		if err != nil {
			return fmt.Errorf("erro ao verificar unicidade do material: %w", err)
		}
		if exists {
			return fmt.Errorf("material já existe com esse ID externo")
		}
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO materiais (titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, pdf_url, disponivel, externo_id, fonte)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			  RETURNING id`
	err = tx.QueryRowContext(ctx, query,
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.ExternoID, m.Fonte,
	).Scan(&m.ID)
	if err != nil {
		return fmt.Errorf("erro ao criar material: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao commitar transação: %w", err)
	}
	return nil
}

func (r *MaterialPostgres) Atualizar(ctx context.Context, m *material.Material) error {
	_, err := r.DB.ExecContext(ctx,
		"UPDATE materiais SET titulo=$1, autor=$2, isbn=$3, categoria=$4, ano_publicacao=$5, descricao=$6, capa_url=$7, pdf_url=$8, disponivel=$9 WHERE id=$10",
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.ID,
	)
	return err
}

func (r *MaterialPostgres) Deletar(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM materiais WHERE id=$1", id)
	return err
}

func (r *MaterialPostgres) SalvarEmprestimo(ctx context.Context, e *material.Emprestimo) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO emprestimos (usuario_id, material_id, data_emprestimo, data_devolucao, status) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		e.UsuarioID, e.MaterialID, e.DataEmprestimo, e.DataDevolucao, e.Status,
	).Scan(&e.ID)
}

func (r *MaterialPostgres) ListarEmprestimosPorUsuario(ctx context.Context, usuarioID int) ([]material.Emprestimo, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, usuario_id, material_id, data_emprestimo, data_devolucao, status FROM emprestimos WHERE usuario_id = $1", usuarioID)
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
func (r *MaterialPostgres) SalvarAvaliacao(ctx context.Context, a *material.Avaliacao) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO avaliacoes (usuario_id, material_id, nota, comentario, data) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		a.UsuarioID, a.MaterialID, a.Nota, a.Comentario, a.Data,
	).Scan(&a.ID)
}

func (r *MaterialPostgres) ListarAvaliacoesPorMaterial(ctx context.Context, materialID int) ([]material.Avaliacao, error) {
	rows, err := r.DB.QueryContext(ctx, "SELECT id, usuario_id, material_id, nota, comentario, data FROM avaliacoes WHERE material_id = $1", materialID)
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

func (r *MaterialPostgres) AdicionarFavorito(ctx context.Context, f *material.Favorito) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO favoritos (usuario_id, material_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", f.UsuarioID, f.MaterialID)
	return err
}

func (r *MaterialPostgres) RemoverFavorito(ctx context.Context, usuarioID, materialID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM favoritos WHERE usuario_id = $1 AND material_id = $2", usuarioID, materialID)
	return err
}

func (r *MaterialPostgres) ListarFavoritosPorUsuario(ctx context.Context, usuarioID int) ([]material.Material, error) {
	query := `
		SELECT m.id, m.titulo, m.autor, COALESCE(m.isbn, ''), m.categoria, m.ano_publicacao,
		       COALESCE(m.descricao, ''), COALESCE(m.capa_url, ''), COALESCE(m.pdf_url, ''),
		       m.disponivel, COALESCE(m.media_nota, 0.0), COALESCE(m.total_avaliacoes, 0)
		FROM materiais m
		JOIN favoritos f ON m.id = f.material_id
		WHERE f.usuario_id = $1`
	rows, err := r.DB.QueryContext(ctx, query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) RegistrarLeitura(ctx context.Context, h *material.HistoricoLeitura) error {
	return r.DB.QueryRowContext(ctx,
		"INSERT INTO historico_leitura (usuario_id, material_id, data) VALUES ($1, $2, $3) RETURNING id",
		h.UsuarioID, h.MaterialID, h.Data,
	).Scan(&h.ID)
}

func (r *MaterialPostgres) ListarHistoricoPorUsuario(ctx context.Context, usuarioID int) ([]material.Material, error) {
	query := `
		SELECT DISTINCT m.id, m.titulo, m.autor, COALESCE(m.isbn, ''), m.categoria, m.ano_publicacao,
		       COALESCE(m.descricao, ''), COALESCE(m.capa_url, ''), COALESCE(m.pdf_url, ''),
		       m.disponivel, COALESCE(m.media_nota, 0.0), COALESCE(m.total_avaliacoes, 0)
		FROM materiais m
		JOIN historico_leitura h ON m.id = h.material_id
		WHERE h.usuario_id = $1
		ORDER BY m.id DESC`
	rows, err := r.DB.QueryContext(ctx, query, usuarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) ObterRecomendacoes(ctx context.Context, usuarioID int, limit int) ([]material.Material, error) {
	// Simple recommendation: match user interests (categories) or just top rated if no interests
	query := `
		SELECT DISTINCT m.id, m.titulo, m.autor, COALESCE(m.isbn, ''), m.categoria, m.ano_publicacao,
		       COALESCE(m.descricao, ''), COALESCE(m.capa_url, ''), COALESCE(m.pdf_url, ''),
		       m.disponivel, COALESCE(m.media_nota, 0.0), COALESCE(m.total_avaliacoes, 0)
		FROM materiais m
		WHERE m.categoria IN (SELECT interesse FROM interesses_usuario WHERE usuario_id = $1)
		OR COALESCE(m.media_nota, 0.0) >= 4.0
		ORDER BY COALESCE(m.media_nota, 0.0) DESC, COALESCE(m.total_avaliacoes, 0) DESC
		LIMIT $2`
	rows, err := r.DB.QueryContext(ctx, query, usuarioID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := scanMaterial(rows, &m); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func scanMaterial(scanner interface {
	Scan(dest ...interface{}) error
}, m *material.Material) error {
	return scanner.Scan(
		&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao,
		&m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes,
	)
}
