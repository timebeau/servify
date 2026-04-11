package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPSessionIPIntelligenceDescribeIP(t *testing.T) {
	var gotAuth string
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotPath = r.URL.Path
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"network_label":"public","location_label":"geo:test-region"}`))
	}))
	defer srv.Close()

	provider := NewHTTPSessionIPIntelligence(srv.URL+"/lookup/{ip}", "secret-token", "Authorization", time.Second)
	desc := provider.DescribeIP("8.8.8.8")

	if gotAuth != "Bearer secret-token" {
		t.Fatalf("expected bearer token header, got %q", gotAuth)
	}
	if gotPath != "/lookup/8.8.8.8" {
		t.Fatalf("expected ip lookup path, got %q", gotPath)
	}
	if desc.NetworkLabel != "public" || desc.LocationLabel != "geo:test-region" {
		t.Fatalf("unexpected description: %+v", desc)
	}
}

func TestHTTPSessionIPIntelligenceDescribeIPReturnsEmptyOnFailure(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	provider := NewHTTPSessionIPIntelligence(srv.URL, "", "", time.Second)
	desc := provider.DescribeIP("8.8.4.4")
	if desc.NetworkLabel != "" || desc.LocationLabel != "" {
		t.Fatalf("expected empty description on upstream failure, got %+v", desc)
	}
}

func TestHTTPSessionIPIntelligenceDescribeIPSupportsNestedDataPayload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"network_label":"public","location_label":"geo:ap-southeast-1"}}`))
	}))
	defer srv.Close()

	provider := NewHTTPSessionIPIntelligence(srv.URL+"/lookup/{ip}", "", "", time.Second)
	desc := provider.DescribeIP("1.1.1.1")

	if desc.NetworkLabel != "public" || desc.LocationLabel != "geo:ap-southeast-1" {
		t.Fatalf("unexpected nested payload description: %+v", desc)
	}
}

func TestHTTPSessionIPIntelligenceDescribeIPUsesCustomAuthHeaderWithoutBearerPrefix(t *testing.T) {
	var gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotHeader = r.Header.Get("X-Geo-Key")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"network_label":"public","location_label":"geo:cn-zj"}`))
	}))
	defer srv.Close()

	provider := NewHTTPSessionIPIntelligence(srv.URL, "geo-secret", "X-Geo-Key", time.Second)
	desc := provider.DescribeIP("8.8.4.4")

	if gotHeader != "geo-secret" {
		t.Fatalf("expected raw custom auth header, got %q", gotHeader)
	}
	if desc.NetworkLabel != "public" || desc.LocationLabel != "geo:cn-zj" {
		t.Fatalf("unexpected description: %+v", desc)
	}
}
