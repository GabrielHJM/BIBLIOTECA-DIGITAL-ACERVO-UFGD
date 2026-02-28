package material

type Material struct {
	ID              int      `json:"id"`
	Titulo          string   `json:"titulo"`
	Autor           string   `json:"autor"`
	ISBN            string   `json:"isbn"`
	Categoria       string   `json:"categoria"`
	AnoPublicacao   int      `json:"ano_publicacao"`
	Paginas         int      `json:"paginas"`
	Descricao       string   `json:"descricao"`
	CapaURL         string   `json:"capa_url"`
	PDFURL          string   `json:"pdf_url"`
	Disponivel      bool     `json:"disponivel"`
	Tags            []string `json:"tags"`
	MediaNota       float64  `json:"media_nota"`
	TotalAvaliacoes int      `json:"total_avaliacoes"`
	ExternoID       string   `json:"externo_id"`
	Fonte           string   `json:"fonte"`
	Status          string   `json:"status"`     // 'pendente', 'aprovado', 'rejeitado'
	CuradorID       int      `json:"curador_id"` // ID do curador que aprovou
}
