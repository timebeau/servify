package dify

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientRetrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("authorization = %q", got)
		}
		if r.URL.Path != "/datasets/ds-1/retrieve" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"query": "refund",
			"records": []map[string]any{
				{
					"segment": map[string]any{
						"id":      "seg-1",
						"content": "Refunds are supported within 7 days.",
					},
					"document": map[string]any{
						"id":   "doc-1",
						"name": "Refund Policy",
					},
					"score": 0.93,
				},
			},
		})
	}))
	defer server.Close()

	client := NewClient(&Config{BaseURL: server.URL, APIKey: "test-key"})
	resp, err := client.Retrieve(context.Background(), "ds-1", &RetrieveRequest{Query: "refund"})
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(resp.Records) != 1 {
		t.Fatalf("records = %d", len(resp.Records))
	}
	if resp.Records[0].Title != "Refund Policy" {
		t.Fatalf("title = %q", resp.Records[0].Title)
	}
}

func TestClientHealthCheck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"id": "ds-1", "name": "Knowledge"})
	}))
	defer server.Close()

	client := NewClient(&Config{BaseURL: server.URL})
	if err := client.HealthCheck(context.Background(), "ds-1"); err != nil {
		t.Fatalf("HealthCheck() error = %v", err)
	}
}
