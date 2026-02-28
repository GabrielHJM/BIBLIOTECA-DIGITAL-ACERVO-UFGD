package handler

import (
	"database/sql"
	"net/http"
	"strconv"
)

func RegisterStatsRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/usuario/estatisticas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			JSONError(w, "Método inválido", http.StatusMethodNotAllowed)
			return
		}

		usuarioID, err := strconv.Atoi(r.URL.Query().Get("usuario_id"))
		if err != nil {
			JSONError(w, "ID de usuário inválido", http.StatusBadRequest)
			return
		}

		// Estatísticas Agregadas e Gamificação
		stats := struct {
			TotalLidos         int            `json:"total_lidos"`
			TotalPaginas       int            `json:"total_paginas"`
			Categorias         map[string]int `json:"categorias"`
			LeiturasPorMes     map[string]int `json:"leituras_por_mes"`
			Badges             []string       `json:"badges"`
			MetaPaginasSemana  int            `json:"meta_paginas_semana"`
			PaginasLidasSemana int            `json:"paginas_lidas_semana"`
		}{
			Categorias:     make(map[string]int),
			LeiturasPorMes: make(map[string]int),
			Badges:         []string{},
		}

		// 0. Meta Semanal do Usuário
		err = db.QueryRowContext(r.Context(), `SELECT COALESCE(meta_paginas_semana, 100) FROM usuarios WHERE id = $1`, usuarioID).Scan(&stats.MetaPaginasSemana)
		if err != nil {
			stats.MetaPaginasSemana = 100 // fallback
		}

		// 1. Total lidos e páginas
		err = db.QueryRowContext(r.Context(), `
			SELECT COUNT(DISTINCT material_id), COALESCE(SUM(m.paginas), 0)
			FROM historico_leitura h
			JOIN materiais m ON h.material_id = m.id
			WHERE h.usuario_id = $1`, usuarioID).Scan(&stats.TotalLidos, &stats.TotalPaginas)
		if err != nil {
			JSONError(w, "Erro ao buscar totais: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// 2. Coletar categorias favoritas
		rowsCat, err := db.QueryContext(r.Context(), `
			SELECT m.categoria, COUNT(*) as qtd
			FROM historico_leitura h
			JOIN materiais m ON h.material_id = m.id
			WHERE h.usuario_id = $1
			GROUP BY m.categoria
			ORDER BY qtd DESC`, usuarioID)
		if err == nil {
			defer rowsCat.Close()
			for rowsCat.Next() {
				var cat string
				var qtd int
				if err := rowsCat.Scan(&cat, &qtd); err == nil {
					stats.Categorias[cat] = qtd
				}
			}
		}

		// 3. Leituras ao longo do tempo (últimos 6 meses)
		rowsTime, err := db.QueryContext(r.Context(), `
			SELECT TO_CHAR(data, 'YYYY-MM') as mes, COUNT(*)
			FROM historico_leitura
			WHERE usuario_id = $1 AND data >= NOW() - INTERVAL '6 months'
			GROUP BY mes
			ORDER BY mes ASC`, usuarioID)
		if err == nil {
			defer rowsTime.Close()
			for rowsTime.Next() {
				var mes string
				var qtd int
				if err := rowsTime.Scan(&mes, &qtd); err == nil {
					stats.LeiturasPorMes[mes] = qtd
				}
			}
		}

		// 4. Progresso Semanal de Leitura (Últimos 7 dias)
		err = db.QueryRowContext(r.Context(), `
			SELECT COALESCE(SUM(m.paginas), 0)
			FROM historico_leitura h
			JOIN materiais m ON h.material_id = m.id
			WHERE h.usuario_id = $1 AND h.data >= NOW() - INTERVAL '7 days'
		`, usuarioID).Scan(&stats.PaginasLidasSemana)
		if err != nil {
			stats.PaginasLidasSemana = 0
		}

		// 5. Lógica de Gamificação (Badges)
		if stats.TotalLidos >= 1 {
			stats.Badges = append(stats.Badges, "Leitor Iniciante")
		}
		if stats.TotalLidos >= 10 {
			stats.Badges = append(stats.Badges, "Rato de Biblioteca")
		}
		if stats.TotalLidos >= 50 {
			stats.Badges = append(stats.Badges, "Mestre Leitor")
		}
		if stats.Categorias["TECNOLOGIA"] >= 5 {
			stats.Badges = append(stats.Badges, "Especialista em Tecnologia")
		}
		if stats.Categorias["SAÚDE"] >= 5 {
			stats.Badges = append(stats.Badges, "Especialista em Saúde")
		}
		if stats.Categorias["MATEMÁTICA"] >= 5 || stats.Categorias["CIÊNCIAS"] >= 5 {
			stats.Badges = append(stats.Badges, "Cientista")
		}

		JSONSuccess(w, stats, http.StatusOK)
	})
}
