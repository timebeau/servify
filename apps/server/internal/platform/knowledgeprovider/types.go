package knowledgeprovider

// SearchRequest is the shared input model for retrieval providers.
type SearchRequest struct {
	Query      string  `json:"query"`
	KnowledgeID string `json:"knowledge_id,omitempty"`
	TopK       int     `json:"top_k,omitempty"`
	Threshold  float64 `json:"threshold,omitempty"`
	Strategy   string  `json:"strategy,omitempty"`
}

// KnowledgeHit is the shared output model for retrieval providers.
type KnowledgeHit struct {
	DocumentID string                 `json:"document_id"`
	Title      string                 `json:"title"`
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
	Source     string                 `json:"source,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
}

// KnowledgeDocument is the shared document input model for indexing providers.
type KnowledgeDocument struct {
	ID       string                 `json:"id,omitempty"`
	Title    string                 `json:"title"`
	Content  string                 `json:"content"`
	Tags     []string               `json:"tags,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
