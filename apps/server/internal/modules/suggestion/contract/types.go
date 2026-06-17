package contract

import "time"

type IntentSuggestion struct {
	Label      string   `json:"label"`
	Confidence float64  `json:"confidence"`
	Matches    []string `json:"matches,omitempty"`
}

type TicketSuggestion struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Category  string    `json:"category"`
	Priority  string    `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	Score     float64   `json:"score"`
}

type KnowledgeDocSuggestion struct {
	ID       uint    `json:"id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Tags     string  `json:"tags"`
	Score    float64 `json:"score"`
}

type SuggestionResponse struct {
	Query          string                   `json:"query"`
	Intent         IntentSuggestion         `json:"intent"`
	SimilarTickets []TicketSuggestion       `json:"similar_tickets"`
	KnowledgeDocs  []KnowledgeDocSuggestion `json:"knowledge_docs"`
	Meta           map[string]interface{}   `json:"meta,omitempty"`
}

type SuggestionRequest struct {
	Query              string `json:"query" binding:"required"`
	TicketLimit        int    `json:"ticket_limit"`
	KnowledgeDocLimit  int    `json:"knowledge_doc_limit"`
	CandidateTicketMax int    `json:"candidate_ticket_max"`
}
