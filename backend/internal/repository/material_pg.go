package repository

import (
	"biblioteca-digital-api/internal/domain/material"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type MaterialPostgres struct {
	DB *sql.DB
}

const materialColumns = `id, titulo, autor, COALESCE(isbn, ''), categoria, ano_publicacao,
                          COALESCE(descricao, ''), COALESCE(capa_url, ''), COALESCE(pdf_url, ''),
                          disponivel, COALESCE(media_nota, 0.0), COALESCE(total_avaliacoes, 0),
                          COALESCE(paginas, 0), COALESCE(externo_id, ''), COALESCE(fonte, ''),
                          COALESCE(status, 'aprovado'), COALESCE(curador_id, 0),
                          COALESCE(dificuldade, 1), COALESCE(xp, 10), COALESCE(relevancia, 0),
                          data_criacao, deleted_at`

func (r *MaterialPostgres) Listar(ctx context.Context, limit, offset int) ([]material.Material, error) {
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE status = 'aprovado' AND deleted_at IS NULL LIMIT $1 OFFSET $2", materialColumns)
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
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE id = $1 AND deleted_at IS NULL", materialColumns)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes, &m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID, &m.Dificuldade, &m.XP, &m.Relevancia, &m.DataCriacao, &m.DeletedAt)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *MaterialPostgres) Pesquisar(ctx context.Context, termo, categoria, fonte string, anoInicio, anoFim int, tags []string, limit, offset int, sort string) ([]material.Material, error) {
	// Base query com Ranking de Relevância FTS
	query := fmt.Sprintf("SELECT %s", materialColumns)

	args := []interface{}{}
	argCount := 1

	if termo != "" {
		// Intelligent FTS: Use websearch_to_tsquery for natural language and support prefix matching (:*)
		// We split the term and add :* to each word to allow "dent" to match "dentista"
		words := strings.Fields(termo)
		ftsQueryParts := []string{}
		for _, w := range words {
			if len(w) > 2 {
				ftsQueryParts = append(ftsQueryParts, w+":*")
			} else {
				ftsQueryParts = append(ftsQueryParts, w)
			}
		}
		formattedQuery := strings.Join(ftsQueryParts, " & ")

		query += fmt.Sprintf(", ts_rank_cd(search_vector, to_tsquery('portuguese', $%d)) as rank", argCount)
		query += " FROM materiais WHERE status = 'aprovado' AND deleted_at IS NULL"
		
		// Hybrid search: FTS with prefix support OR ILIKE fallback for unaccented matches
		query += fmt.Sprintf(" AND (search_vector @@ to_tsquery('portuguese', $%d) OR unaccent(titulo) ILIKE unaccent($%d) OR unaccent(autor) ILIKE unaccent($%d))", argCount, argCount+1, argCount+1)

		args = append(args, formattedQuery, "%"+termo+"%")
		argCount += 2
	} else {
		query += " FROM materiais WHERE status = 'aprovado' AND deleted_at IS NULL"
	}

	if categoria != "" {
		query += fmt.Sprintf(" AND categoria ILIKE $%d", argCount)
		args = append(args, "%"+categoria+"%")
		argCount++
	}

	if fonte != "" {
		query += fmt.Sprintf(" AND unaccent(fonte) ILIKE unaccent($%d)", argCount)
		args = append(args, "%"+fonte+"%")
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

	// Ordenação
	orderBy := "unaccent(LOWER(titulo)) ASC"

	switch sort {
	case "az":
		orderBy = "unaccent(LOWER(titulo)) ASC"
	case "za":
		orderBy = "unaccent(LOWER(titulo)) DESC"
	case "recent":
		orderBy = "ano_publicacao DESC, id DESC"
	case "oldest":
		orderBy = "ano_publicacao ASC, id DESC"
	case "random":
		orderBy = "RANDOM()"
	case "relevancia", "":
		if termo != "" {
			orderBy = "rank DESC, relevancia DESC, id DESC"
		} else {
			orderBy = "relevancia DESC, total_avaliacoes DESC, id DESC"
		}
	}

	query += " ORDER BY " + orderBy
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("erro na busca FTS: %w", err)
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		var rank float64
		dest := []interface{}{
			&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao,
			&m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes,
			&m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID,
			&m.Dificuldade, &m.XP, &m.Relevancia, &m.DataCriacao, &m.DeletedAt,
		}
		if termo != "" {
			dest = append(dest, &rank)
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) Sugerir(ctx context.Context, termo string, limit int) ([]material.Material, error) {
	// Query otimizada para sugestão (apenas campos necessários para UI de dropdown)
	query := `SELECT id, titulo, autor, categoria, COALESCE(capa_url, ''), media_nota
	          FROM materiais
	          WHERE status = 'aprovado' AND deleted_at IS NULL
	          AND (unaccent(titulo) ILIKE unaccent($1) OR unaccent(autor) ILIKE unaccent($1))
	          ORDER BY relevancia DESC, media_nota DESC
	          LIMIT $2`

	rows, err := r.DB.QueryContext(ctx, query, "%"+termo+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materiais []material.Material
	for rows.Next() {
		var m material.Material
		if err := rows.Scan(&m.ID, &m.Titulo, &m.Autor, &m.Categoria, &m.CapaURL, &m.MediaNota); err != nil {
			return nil, err
		}
		materiais = append(materiais, m)
	}
	return materiais, nil
}

func (r *MaterialPostgres) BuscarSimilares(ctx context.Context, materialID int, limit int) ([]material.Material, error) {
	query := fmt.Sprintf(`SELECT %s
	          FROM materiais m1
	          JOIN materiais m2 ON m1.categoria = m2.categoria AND m1.id != m2.id
	          WHERE m1.id = $1 AND m2.status = 'aprovado' AND m2.deleted_at IS NULL
	          LIMIT $2`, "m2."+replaceColumns(materialColumns, "m2."))
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
	query := `INSERT INTO materiais (titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, pdf_url, disponivel, externo_id, fonte, paginas, dificuldade, xp, relevancia, status)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, 'aprovado')
			  ON CONFLICT (externo_id) WHERE externo_id IS NOT NULL DO UPDATE SET
			    titulo = EXCLUDED.titulo,
				autor = EXCLUDED.autor,
				capa_url = COALESCE(EXCLUDED.capa_url, materiais.capa_url),
				pdf_url = COALESCE(EXCLUDED.pdf_url, materiais.pdf_url),
				status = 'aprovado'
			  RETURNING id`

	err := r.DB.QueryRowContext(ctx, query,
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.ExternoID, m.Fonte, m.Paginas, m.Dificuldade, m.XP, m.Relevancia,
	).Scan(&m.ID)

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "já existe") {
			// Se falhou por conflito de ISBN ou Titulo/Autor (que não estão no ON CONFLICT acima), tentamos buscar o ID
			var id int
			if m.ISBN != "" {
				_ = r.DB.QueryRowContext(ctx, "SELECT id FROM materiais WHERE isbn = $1", m.ISBN).Scan(&id)
			} else {
				_ = r.DB.QueryRowContext(ctx, "SELECT id FROM materiais WHERE titulo = $1 AND autor = $2", m.Titulo, m.Autor).Scan(&id)
			}
			if id > 0 {
				m.ID = id
				return nil
			}
		}
		return fmt.Errorf("erro ao criar material: %w", err)
	}
	return nil
}

func (r *MaterialPostgres) Atualizar(ctx context.Context, m *material.Material) error {
	_, err := r.DB.ExecContext(ctx,
		"UPDATE materiais SET titulo=$1, autor=$2, isbn=$3, categoria=$4, ano_publicacao=$5, descricao=$6, capa_url=$7, pdf_url=$8, disponivel=$9, paginas=$10 WHERE id=$11",
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.Paginas, m.ID,
	)
	return err
}

func (r *MaterialPostgres) Deletar(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, "UPDATE materiais SET deleted_at = CURRENT_TIMESTAMP WHERE id=$1", id)
	return err
}

func (r *MaterialPostgres) AdicionarFavorito(ctx context.Context, f *material.Favorito) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO favoritos (usuario_id, material_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", f.UsuarioID, f.MaterialID)
	return err
}

func (r *MaterialPostgres) RemoverFavorito(ctx context.Context, usuarioID, materialID int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM favoritos WHERE usuario_id = $1 AND material_id = $2", usuarioID, materialID)
	return err
}

func (r *MaterialPostgres) Avaliar(ctx context.Context, usuarioID, materialID int, nota float64) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Inserir ou atualizar a avaliação individual
	_, err = tx.ExecContext(ctx, `
		INSERT INTO avaliacoes (usuario_id, material_id, nota)
		VALUES ($1, $2, $3)
		ON CONFLICT (usuario_id, material_id) DO UPDATE SET nota = EXCLUDED.nota`,
		usuarioID, materialID, nota)
	if err != nil {
		return err
	}

	// 2. Recalcular a média e o total para o material
	var media float64
	var total int
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(AVG(nota), 0), COUNT(*)
		FROM avaliacoes
		WHERE material_id = $1`, materialID).Scan(&media, &total)
	if err != nil {
		return err
	}

	// 3. Atualizar a tabela de materiais
	_, err = tx.ExecContext(ctx, `
		UPDATE materiais
		SET media_nota = $1, total_avaliacoes = $2
		WHERE id = $3`, media, total, materialID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *MaterialPostgres) ListarFavoritosPorUsuario(ctx context.Context, usuarioID int) ([]material.Material, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM materiais m
		JOIN favoritos f ON m.id = f.material_id
		WHERE f.usuario_id = $1 AND m.deleted_at IS NULL`, replaceColumns(materialColumns, "m."))

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
	query := fmt.Sprintf(`
		SELECT %s FROM materiais m
		JOIN (
			SELECT material_id, MAX(data) as last_read
			FROM historico_leitura
			WHERE usuario_id = $1
			GROUP BY material_id
		) h ON m.id = h.material_id
		WHERE m.deleted_at IS NULL
		ORDER BY h.last_read DESC`, replaceColumns(materialColumns, "m."))
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

func scanMaterial(scanner interface {
	Scan(dest ...interface{}) error
}, m *material.Material) error {
	return scanner.Scan(
		&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao,
		&m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes,
		&m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID, &m.Dificuldade, &m.XP, &m.Relevancia,
		&m.DataCriacao, &m.DeletedAt,
	)
}

func replaceColumns(cols, alias string) string {
	parts := strings.Split(cols, ",")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "COALESCE(") {
			// replace COALESCE(field, ...) with COALESCE(alias.field, ...)
			parts[i] = strings.Replace(p, "COALESCE(", "COALESCE("+alias, 1)
		} else if p == "" || p[0] == '\'' || (p[0] >= '0' && p[0] <= '9') {
			// It's a literal or closing part of COALESCE, don't alias
			parts[i] = p
		} else {
			parts[i] = alias + p
		}
	}
	return strings.Join(parts, ", ")
}

func (r *MaterialPostgres) ListarPendentes(ctx context.Context) ([]material.Material, error) {
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE status = 'pendente' AND deleted_at IS NULL ORDER BY id DESC", materialColumns)
	rows, err := r.DB.QueryContext(ctx, query)
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

func (r *MaterialPostgres) AtualizarStatus(ctx context.Context, id int, status string, curadorID int) error {
	query := `UPDATE materiais SET status = $1, curador_id = $2 WHERE id = $3`
	_, err := r.DB.ExecContext(ctx, query, status, curadorID, id)
	return err
}

func (r *MaterialPostgres) ObterMetricasGlobais(ctx context.Context) (map[string]interface{}, error) {
	var totalUsuarios, totalMateriais, totalLeituras int

	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM usuarios").Scan(&totalUsuarios); err != nil {
		return nil, fmt.Errorf("erro ao contar usuários: %w", err)
	}
	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM materiais").Scan(&totalMateriais); err != nil {
		return nil, fmt.Errorf("erro ao contar materiais: %w", err)
	}
	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM historico_leitura").Scan(&totalLeituras); err != nil {
		return nil, fmt.Errorf("erro ao contar leituras: %w", err)
	}

	rows, err := r.DB.QueryContext(ctx, "SELECT fonte, COUNT(*) FROM materiais GROUP BY fonte")
	if err != nil {
		return nil, fmt.Errorf("erro ao agrupar fontes: %w", err)
	}
	defer rows.Close()

	fontes := make(map[string]int)
	for rows.Next() {
		var fonte string
		var count int
		if err := rows.Scan(&fonte, &count); err == nil {
			if fonte == "" {
				fonte = "Local"
			}
			fontes[fonte] = count
		}
	}

	return map[string]interface{}{
		"total_usuarios":  totalUsuarios,
		"total_materiais": totalMateriais,
		"total_leituras":  totalLeituras,
		"fontes":          fontes,
	}, nil
}
