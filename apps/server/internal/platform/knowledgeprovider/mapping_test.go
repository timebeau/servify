package knowledgeprovider

import "testing"

func TestResolveNamespace(t *testing.T) {
	ns := ResolveNamespace("tenant-default", "kb-default", "", "")
	if ns.TenantID != "tenant-default" || ns.KnowledgeID != "kb-default" {
		t.Fatalf("unexpected namespace: %+v", ns)
	}

	ns = ResolveNamespace("tenant-default", "", "tenant-a", "")
	if ns.TenantID != "tenant-a" || ns.KnowledgeID != "tenant-a" {
		t.Fatalf("expected tenant fallback to knowledge id, got %+v", ns)
	}
}
