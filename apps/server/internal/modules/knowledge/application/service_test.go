package application

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"servify/apps/server/internal/modules/knowledge/domain"
	"servify/apps/server/internal/platform/knowledgeprovider"
	memorykp "servify/apps/server/internal/platform/knowledgeprovider/memory"
	mockkp "servify/apps/server/internal/platform/knowledgeprovider/mock"
)

type memDocRepo struct {
	docs map[string]*domain.Document
}

func (r *memDocRepo) Create(ctx context.Context, doc *domain.Document) error {
	if r.docs == nil {
		r.docs = map[string]*domain.Document{}
	}
	if doc.ID == "" {
		doc.ID = fmt.Sprintf("doc-%d", len(r.docs)+1)
	}
	cp := *doc
	r.docs[doc.ID] = &cp
	return nil
}
func (r *memDocRepo) Update(ctx context.Context, doc *domain.Document) error {
	cp := *doc
	r.docs[doc.ID] = &cp
	return nil
}
func (r *memDocRepo) Delete(ctx context.Context, id string) error {
	delete(r.docs, id)
	return nil
}
func (r *memDocRepo) Get(ctx context.Context, id string) (*domain.Document, error) {
	doc, ok := r.docs[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	cp := *doc
	return &cp, nil
}
func (r *memDocRepo) List(ctx context.Context, filter ListDocumentsFilter) ([]domain.Document, int64, error) {
	out := make([]domain.Document, 0, len(r.docs))
	for _, doc := range r.docs {
		if filter.Category != "" && doc.Category != filter.Category {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(doc.Title+doc.Content), strings.ToLower(filter.Search)) {
			continue
		}
		out = append(out, *doc)
	}
	return out, int64(len(out)), nil
}

type memJobRepo struct {
	jobs map[string]*domain.IndexJob
}

func (r *memJobRepo) Create(ctx context.Context, job *domain.IndexJob) error {
	if r.jobs == nil {
		r.jobs = map[string]*domain.IndexJob{}
	}
	if job.ID == "" {
		job.ID = fmt.Sprintf("job-%d", len(r.jobs)+1)
	}
	cp := *job
	r.jobs[job.ID] = &cp
	return nil
}
func (r *memJobRepo) Update(ctx context.Context, job *domain.IndexJob) error {
	cp := *job
	r.jobs[job.ID] = &cp
	return nil
}
func (r *memJobRepo) Get(ctx context.Context, id string) (*domain.IndexJob, error) {
	job, ok := r.jobs[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	cp := *job
	return &cp, nil
}

func TestServiceCreateDocument(t *testing.T) {
	svc := NewService(&memDocRepo{}, &memJobRepo{}, nil)
	doc, err := svc.CreateDocument(context.Background(), CreateDocumentRequest{
		Title:   " KB Title ",
		Content: " KB Content ",
		Tags:    []string{" a ", "", "b"},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if doc.ID == "" {
		t.Fatal("expected generated id")
	}
	if len(doc.Tags) != 2 {
		t.Fatalf("unexpected tags: %+v", doc.Tags)
	}
}

func TestServiceRunIndexJob(t *testing.T) {
	docRepo := &memDocRepo{}
	jobRepo := &memJobRepo{}
	provider := &mockkp.Provider{}
	svc := NewService(docRepo, jobRepo, provider)

	doc, err := svc.CreateDocument(context.Background(), CreateDocumentRequest{
		ID:      "doc-1",
		Title:   "Billing",
		Content: "Billing details",
	})
	if err != nil {
		t.Fatalf("create doc: %v", err)
	}
	job, err := svc.QueueIndexJob(context.Background(), QueueIndexJobRequest{
		JobID:      "job-1",
		DocumentID: doc.ID,
	})
	if err != nil {
		t.Fatalf("queue job: %v", err)
	}
	result, err := svc.RunIndexJob(context.Background(), RunIndexJobRequest{JobID: job.ID})
	if err != nil {
		t.Fatalf("run job: %v", err)
	}
	if result.Status != string(domain.IndexJobDone) {
		t.Fatalf("unexpected status: %+v", result)
	}
	if provider.Documents["doc-1"].Title != "Billing" {
		t.Fatalf("expected provider upsert, got %+v", provider.Documents)
	}
}

func TestServiceListDocuments(t *testing.T) {
	docRepo := &memDocRepo{}
	svc := NewService(docRepo, &memJobRepo{}, nil)
	_, _ = svc.CreateDocument(context.Background(), CreateDocumentRequest{ID: "doc-1", Title: "Billing", Content: "billing content", Category: "faq"})
	_, _ = svc.CreateDocument(context.Background(), CreateDocumentRequest{ID: "doc-2", Title: "Support", Content: "support content", Category: "guide"})

	docs, total, err := svc.ListDocuments(context.Background(), ListDocumentsFilter{
		Category: "faq",
	})
	if err != nil {
		t.Fatalf("list docs: %v", err)
	}
	if total != 1 || len(docs) != 1 {
		t.Fatalf("unexpected docs: total=%d docs=%d", total, len(docs))
	}
}

func TestServiceRunIndexJobProviderSwitchRegression(t *testing.T) {
	providers := map[string]knowledgeprovider.KnowledgeProvider{
		"mock":   &mockkp.Provider{},
		"memory": memorykp.NewProvider("tenant-a", "kb-a"),
	}

	for name, provider := range providers {
		t.Run(name, func(t *testing.T) {
			docRepo := &memDocRepo{}
			jobRepo := &memJobRepo{}
			svc := NewService(docRepo, jobRepo, provider)

			doc, err := svc.CreateDocument(context.Background(), CreateDocumentRequest{
				ID:      "doc-switch",
				Title:   "Billing",
				Content: "Billing details",
			})
			if err != nil {
				t.Fatalf("create doc: %v", err)
			}
			job, err := svc.QueueIndexJob(context.Background(), QueueIndexJobRequest{
				JobID:      "job-switch",
				DocumentID: doc.ID,
			})
			if err != nil {
				t.Fatalf("queue job: %v", err)
			}
			result, err := svc.RunIndexJob(context.Background(), RunIndexJobRequest{JobID: job.ID})
			if err != nil {
				t.Fatalf("run job: %v", err)
			}
			if result.Status != string(domain.IndexJobDone) {
				t.Fatalf("unexpected status: %+v", result)
			}
		})
	}
}

var _ knowledgeprovider.KnowledgeProvider = (*mockkp.Provider)(nil)
