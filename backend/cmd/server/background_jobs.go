package main

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
	"sync"
	"time"

	"biblioteca-digital-api/internal/harvester"
	"biblioteca-digital-api/internal/pkg/logger"
	"biblioteca-digital-api/internal/repository"

	"go.uber.org/zap"
)

var (
	syncStateMu  sync.Mutex
	categoryPage map[string]int
)

func init() {
	categoryPage = make(map[string]int)
}

func startBackgroundSync(db *sql.DB) {
	repo := &repository.MaterialPostgres{DB: db}
	mh := harvester.NewMultiSourceHarvester()

	// Matriz de Internacionalização (i18n)
	categories := []string{
		// Português
		"TECNOLOGIA", "SAÚDE", "DIREITO", "CIÊNCIAS", "MATEMÁTICA", "EDUCAÇÃO",
		"PROGRAMAÇÃO", "MEDICINA", "ENGENHARIA", "PSICOLOGIA", "BIOLOGIA",
		// Inglês
		"TECHNOLOGY", "HEALTH", "LAW", "SCIENCE", "MATHEMATICS", "EDUCATION",
		"PROGRAMMING", "MEDICINE", "ENGINEERING", "PSYCHOLOGY", "BIOLOGY",
		// Espanhol
		"TECNOLOGÍA", "SALUD", "DERECHO", "CIENCIA", "MATEMÁTICAS", "EDUCACIÓN",
		"PROGRAMACIÓN", "MEDICINA", "INGENIERÍA", "PSICOLOGÍA", "BIOLOGÍA",
		// Francês
		"TECHNOLOGIE", "SANTÉ", "DROIT", "SCIENCE", "MATHÉMATIQUES", "ÉDUCATION",
		"PROGRAMMATION", "MÉDECINE", "INGÉNIERIE", "PSYCHOLOGIE", "BIOLOGIE",
		// Alemão
		"TECHNOLOGIE", "GESUNDHEIT", "RECHT", "WISSENSCHAFT", "MATHEMATIK", "BILDUNG",
		"PROGRAMMIERUNG", "MEDIZIN", "INGENIEURWESEN", "PSYCHOLOGIE", "BIOLOGIE",
	}

	logger.Info("Starting Deep Background Harvester (i18n)", zap.Int("total_categories", len(categories)))

	// Inicializa o state com offset 0 para todas as categorias
	for _, cat := range categories {
		categoryPage[cat] = 0
	}

	// Start the API Supervisor
	harvester.GlobalSupervisor.StartMonitoring(1 * time.Minute)

	// Executa num loop infinito contínuo, com pausas para evitar Rate Limits
	go func() {
		// Pausa inicial para o servidor terminar de subir
		time.Sleep(10 * time.Second)

		for {
			for _, cat := range categories {
				syncDeepBooks(repo, mh, cat)
				// Pausa de 30 segundos entre requisições de categorias
				time.Sleep(30 * time.Second)
			}
			// Quando termina uma rodada completa, espera um pouco mais
			logger.Info("Deep Background Harvester completed a full i18n cycle. Resting...")
			time.Sleep(5 * time.Minute)
		}
	}()
}

func syncDeepBooks(repo *repository.MaterialPostgres, mh *harvester.MultiSourceHarvester, category string) {
	if !syncMu.TryLock() {
		return
	}
	defer syncMu.Unlock()

	syncStateMu.Lock()
	currentOffset := categoryPage[category]
	categoryPage[category] = currentOffset + 20
	syncStateMu.Unlock()

	logger.Info("DeepSync: Harvesting category", zap.String("category", category), zap.Int("offset", currentOffset))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	mats, err := mh.Search(ctx, "", category, "", 0, 0, 20, currentOffset)
	if err != nil {
		logger.Error("DeepSync: Harvester search failed", zap.String("category", category), zap.Error(err))
		return
	}

	count := 0
	for i := range mats {
		if err := repo.Criar(context.Background(), &mats[i]); err != nil {
			if !strings.Contains(err.Error(), "já existe") {
				logger.Debug("DeepSync: Failed to save material", zap.String("title", mats[i].Titulo), zap.Error(err))
			}
		} else {
			count++
		}
	}
	logger.Info("DeepSync: Category populated", zap.String("category", category), zap.Int("new_books_saved", count), zap.Int("offset", currentOffset))
}

// startDeadLinkChecker varre constantemente os materiais disponíveis e testa as URLs dos PDFs
func startDeadLinkChecker(db *sql.DB) {
	ticker := time.NewTicker(10 * time.Second) // Roda a cada 10 segundos
	defer ticker.Stop()

	// Cria um cliente HTTP com timeout curto para não travar
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	logger.Info("Starting Dead-Link Checker")

	for range ticker.C {
		// Busca 1 material que não é checado há muito tempo (ou nunca foi)
		var id int
		var pdfURL string

		query := `SELECT id, pdf_url FROM materiais 
				  WHERE disponivel = true AND status = 'aprovado' AND deleted_at IS NULL
				  ORDER BY last_link_check ASC NULLS FIRST LIMIT 1`
				  
		err := db.QueryRow(query).Scan(&id, &pdfURL)
		if err != nil {
			if err != sql.ErrNoRows {
				logger.Error("DeadLinkChecker: Error fetching material", zap.Error(err))
			}
			continue
		}

		// Testa a URL
		isAlive := checkLinkIsAlive(client, pdfURL)

		// Atualiza o banco de dados
		var updateQuery string
		if isAlive {
			updateQuery = `UPDATE materiais SET last_link_check = CURRENT_TIMESTAMP WHERE id = $1`
			_, err = db.Exec(updateQuery, id)
			if err != nil {
				logger.Error("DeadLinkChecker: Error updating last_link_check", zap.Error(err))
			} else {
				logger.Debug("DeadLinkChecker: Link OK", zap.Int("id", id), zap.String("url", pdfURL))
			}
		} else {
			// Link morto -> Inativa o material
			updateQuery = `UPDATE materiais SET disponivel = false, last_link_check = CURRENT_TIMESTAMP WHERE id = $1`
			_, err = db.Exec(updateQuery, id)
			if err != nil {
				logger.Error("DeadLinkChecker: Error invalidating material", zap.Error(err))
			} else {
				logger.Warn("DeadLinkChecker: Link DEAD. Material deactivated.", zap.Int("id", id), zap.String("url", pdfURL))
			}
		}
	}
}

func checkLinkIsAlive(client *http.Client, urlStr string) bool {
	// Algumas URLs dão pau com HEAD, tentaremos GET range se HEAD falhar ou podemos tentar apenas GET.
	// Vamos usar GET com Range para baixar poucos bytes.
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Range", "bytes=0-100")
	req.Header.Set("User-Agent", "AcervusBot/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 200, 206, ou até 403 (alguns sites dão 403 para robôs mas o link tá vivo)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 || resp.StatusCode == 403 {
		return true
	}

	return false
}
