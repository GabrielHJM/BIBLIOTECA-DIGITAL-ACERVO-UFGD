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
                          COALESCE(dificuldade, 1), COALESCE(xp, 10), COALESCE(relevancia, 0)`

func (r *MaterialPostgres) Listar(ctx context.Context, limit, offset int) ([]material.Material, error) {
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE status = 'aprovado' AND pdf_url ILIKE '%%.pdf%%' AND capa_url NOT LIKE 'https://images.unsplash.com/%%' LIMIT $1 OFFSET $2", materialColumns)
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
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE id = $1", materialColumns)
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes, &m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID, &m.Dificuldade, &m.XP, &m.Relevancia)
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
	ftsQuery := ""

	if termo != "" {
		// Busca FTS em Português
		// Nós vamos fornecer o mesmo 'termo' mais tarde no código principal para FTS, mas aqui registramos que estamos usando o rank()
		// Usamos COALESCE pois se search_vector estiver NULL, ts_rank retorna NULL e quebra o Scan em float64
		query += fmt.Sprintf(", COALESCE(ts_rank(search_vector, plainto_tsquery('portuguese', $%d)), 0.0) as rank", argCount)
		// Não adicionamos args aqui ainda porque reestruturamos a injeção condicional no bloco AND abaixo para unificar.
	}

	query += " FROM materiais WHERE status = 'aprovado' AND pdf_url ILIKE '%%.pdf%%' AND capa_url NOT LIKE 'https://images.unsplash.com/%%'"

	// Busca híbrida: FTS Rankado OR Fallback ILIKE em título/autor
	if termo != "" {
		ftsQuery = fmt.Sprintf("plainto_tsquery('portuguese', $%d)", argCount)
		// Alterado de apenas search_vector @@ ftsQuery para também cobrir casos onde FTS falha
		query += fmt.Sprintf(" AND (search_vector @@ %s OR unaccent(titulo) ILIKE unaccent($%d) OR unaccent(autor) ILIKE unaccent($%d))", ftsQuery, argCount+1, argCount+1)

		// O Select já tem o rank, precisamos adicionar o termo com % para o ILIKE
		args = append(args, termo, "%"+termo+"%")
		argCount += 2 // Adicionamos 2 args (termo original para FTS e termo com % para ILIKE)
	}

	if categoria != "" {
		query += fmt.Sprintf(" AND unaccent(categoria) ILIKE unaccent($%d)", argCount)
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

	// Ordenação Inovadora: Relevância FTS -> Nota -> ID
	orderBy := "id DESC"
	if termo != "" {
		orderBy = "rank DESC, media_nota DESC"
	} else if sort == "random" {
		orderBy = "RANDOM()"
	} else if sort == "rating" {
		orderBy = "media_nota DESC, total_avaliacoes DESC"
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
			&m.Dificuldade, &m.XP, &m.Relevancia,
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

func (r *MaterialPostgres) BuscarSimilares(ctx context.Context, materialID int, limit int) ([]material.Material, error) {
	query := fmt.Sprintf(`SELECT %s
	          FROM materiais m1
	          JOIN materiais m2 ON m1.categoria = m2.categoria AND m1.id != m2.id
	          WHERE m1.id = $1 AND m2.status = 'aprovado' AND m2.pdf_url ILIKE '%%.pdf%%' AND m2.capa_url NOT LIKE 'https://images.unsplash.com/%%'
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
	// O sistema verifica duplicatas no app e já preenche o ID se existir para o cache não quebrar
	if m.ExternoID != "" {
		err := r.DB.QueryRowContext(ctx, "SELECT id FROM materiais WHERE externo_id = $1", m.ExternoID).Scan(&m.ID)
		if err == nil {
			return fmt.Errorf("material já existe com esse ID externo")
		}
	} else {
		if m.ISBN != "" {
			err := r.DB.QueryRowContext(ctx, "SELECT id FROM materiais WHERE isbn = $1 limit 1", m.ISBN).Scan(&m.ID)
			if err == nil {
				return fmt.Errorf("material já existe por isbn")
			}
		} else {
			err := r.DB.QueryRowContext(ctx, "SELECT id FROM materiais WHERE titulo = $1 AND autor = $2 limit 1", m.Titulo, m.Autor).Scan(&m.ID)
			if err == nil {
				return fmt.Errorf("material já existe por titulo e autor")
			}
		}
	}

	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO materiais (titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, pdf_url, disponivel, externo_id, fonte, paginas, dificuldade, xp, relevancia)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			  RETURNING id`
	err = tx.QueryRowContext(ctx, query,
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.ExternoID, m.Fonte, m.Paginas, m.Dificuldade, m.XP, m.Relevancia,
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
		"UPDATE materiais SET titulo=$1, autor=$2, isbn=$3, categoria=$4, ano_publicacao=$5, descricao=$6, capa_url=$7, pdf_url=$8, disponivel=$9, paginas=$10 WHERE id=$11",
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.Paginas, m.ID,
	)
	return err
}

func (r *MaterialPostgres) Deletar(ctx context.Context, id int) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM materiais WHERE id=$1", id)
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

func (r *MaterialPostgres) ListarFavoritosPorUsuario(ctx context.Context, usuarioID int) ([]material.Material, error) {
	query := fmt.Sprintf(`
		SELECT %s
		FROM materiais m
		JOIN favoritos f ON m.id = f.material_id
		WHERE f.usuario_id = $1`, replaceColumns(materialColumns, "m."))
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
		SELECT DISTINCT %s
		FROM materiais m
		JOIN historico_leitura h ON m.id = h.material_id
		WHERE h.usuario_id = $1
		ORDER BY m.id DESC`, replaceColumns(materialColumns, "m."))
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
	)
}

func replaceColumns(cols, alias string) string {
	parts := strings.Split(cols, ",")
	for i, p := range parts {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, "COALESCE(") {
			// replace COALESCE(field, ...) with COALESCE(alias.field, ...)
			parts[i] = strings.Replace(p, "COALESCE(", "COALESCE("+alias, 1)
		} else {
			parts[i] = alias + p
		}
	}
	return strings.Join(parts, ", ")
}

func (r *MaterialPostgres) ListarPendentes(ctx context.Context) ([]material.Material, error) {
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE status = 'pendente' ORDER BY id DESC", materialColumns)
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

	_ = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM usuarios").Scan(&totalUsuarios)
	_ = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM materiais").Scan(&totalMateriais)
	_ = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM historico_leitura").Scan(&totalLeituras)

	rows, _ := r.DB.QueryContext(ctx, "SELECT fonte, COUNT(*) FROM materiais GROUP BY fonte")
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
