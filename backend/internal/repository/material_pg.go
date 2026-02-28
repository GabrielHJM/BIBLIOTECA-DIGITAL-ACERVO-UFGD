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
                          COALESCE(status, 'aprovado'), COALESCE(curador_id, 0)`

func (r *MaterialPostgres) Listar(ctx context.Context, limit, offset int) ([]material.Material, error) {
	query := fmt.Sprintf("SELECT %s FROM materiais WHERE status = 'aprovado' LIMIT $1 OFFSET $2", materialColumns)
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
		&m.ID, &m.Titulo, &m.Autor, &m.ISBN, &m.Categoria, &m.AnoPublicacao, &m.Descricao, &m.CapaURL, &m.PDFURL, &m.Disponivel, &m.MediaNota, &m.TotalAvaliacoes, &m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID)
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
		ftsQuery = fmt.Sprintf("plainto_tsquery('portuguese', $%d)", argCount)
		query += fmt.Sprintf(", ts_rank(search_vector, %s) as rank", ftsQuery)
		args = append(args, termo)
		argCount++
	}

	query += " FROM materiais WHERE status = 'aprovado'"

	if termo != "" {
		query += fmt.Sprintf(" AND search_vector @@ %s", ftsQuery)
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
	          WHERE m1.id = $1 AND m2.status = 'aprovado'
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

	query := `INSERT INTO materiais (titulo, autor, isbn, categoria, ano_publicacao, descricao, capa_url, pdf_url, disponivel, externo_id, fonte, paginas)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
			  RETURNING id`
	err = tx.QueryRowContext(ctx, query,
		m.Titulo, m.Autor, m.ISBN, m.Categoria, m.AnoPublicacao, m.Descricao, m.CapaURL, m.PDFURL, m.Disponivel, m.ExternoID, m.Fonte, m.Paginas,
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

func (r *MaterialPostgres) ObterRecomendacoes(ctx context.Context, usuarioID int, limit int) ([]material.Material, error) {
	// Recomendação Inteligente:
	// 1. Materiais lidos por usuários que compartilham dos mesmos interesses
	// 2. Materiais da mesma categoria de interesse do usuário
	// 3. Materiais mais bem avaliados globalmente
	query := fmt.Sprintf(`
		WITH interesses AS (
			SELECT interesse FROM interesses_usuario WHERE usuario_id = $1
		),
		usuarios_similares AS (
			SELECT DISTINCT usuario_id FROM interesses_usuario
			WHERE interesse IN (SELECT interesse FROM interesses) AND usuario_id != $1
		),
		recomendados AS (
			-- Opção 1: Lidos por usuários similares
			SELECT material_id, 3 as peso FROM historico_leitura WHERE usuario_id IN (SELECT usuario_id FROM usuarios_similares)
			UNION ALL
			-- Opção 2: Mesma categoria de interesse
			SELECT id as material_id, 2 as peso FROM materiais WHERE categoria IN (SELECT interesse FROM interesses)
			UNION ALL
			-- Opção 3: Nota alta (Global)
			SELECT id as material_id, 1 as peso FROM materiais WHERE media_nota >= 4.0
		)
		SELECT %s
		FROM materiais m
		JOIN (
			SELECT material_id, SUM(peso) as relevancia
			FROM recomendados
			GROUP BY material_id
		) r_final ON m.id = r_final.material_id
		WHERE m.id NOT IN (SELECT material_id FROM historico_leitura WHERE usuario_id = $1)
		ORDER BY r_final.relevancia DESC, m.media_nota DESC
		LIMIT $2`, replaceColumns(materialColumns, "m."))
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
		&m.Paginas, &m.ExternoID, &m.Fonte, &m.Status, &m.CuradorID,
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
