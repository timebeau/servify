package dify

type Dataset struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RetrievalModel struct {
	SearchMethod    string  `json:"search_method,omitempty"`
	RerankingEnable bool    `json:"reranking_enable,omitempty"`
	TopK            int     `json:"top_k,omitempty"`
	ScoreThreshold  float64 `json:"score_threshold,omitempty"`
}

type RetrieveRequest struct {
	Query          string         `json:"query"`
	RetrievalModel RetrievalModel `json:"retrieval_model,omitempty"`
}

type RetrieveRecord struct {
	SegmentID  string                 `json:"segment_id"`
	DocumentID string                 `json:"document_id"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

type RetrieveResponse struct {
	Query   string           `json:"query"`
	Records []RetrieveRecord `json:"records"`
}

type ProcessRule struct {
	Mode string `json:"mode,omitempty"`
}

type CreateDocumentRequest struct {
	Name              string         `json:"name"`
	Text              string         `json:"text"`
	IndexingTechnique string         `json:"indexing_technique,omitempty"`
	ProcessRule       ProcessRule    `json:"process_rule,omitempty"`
	RetrievalModel    RetrievalModel `json:"retrieval_model,omitempty"`
}

type Document struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
