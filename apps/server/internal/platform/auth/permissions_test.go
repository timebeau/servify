package auth

import (
	"testing"

	"servify/apps/server/internal/config"
)

func TestHasPermission_WildcardsAndExact(t *testing.T) {
	tests := []struct {
		name     string
		granted  []string
		required string
		want     bool
	}{
		{"star", []string{"*"}, "tickets.read", true},
		{"exact", []string{"tickets.read"}, "tickets.read", true},
		{"prefixStar", []string{"tickets.*"}, "tickets.read", true},
		{"prefixStarNested", []string{"tickets.*"}, "tickets.write", true},
		{"noMatch", []string{"customers.read"}, "tickets.read", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasPermission(tt.granted, tt.required); got != tt.want {
				t.Fatalf("HasPermission(%v, %q)=%v want %v", tt.granted, tt.required, got, tt.want)
			}
		})
	}
}

func TestResolverExpandPermissions(t *testing.T) {
	t.Run("rbac enabled", func(t *testing.T) {
		got := (Resolver{
			RBAC: config.RBACConfig{
				Enabled: true,
				Roles: map[string][]string{
					"viewer": {"tickets.read"},
				},
			},
		}).ExpandPermissions([]string{"viewer"}, []string{"customers.read"})
		if len(got) != 2 || got[0] != "customers.read" || got[1] != "tickets.read" {
			t.Fatalf("unexpected permissions: %v", got)
		}
	})

	t.Run("fallback role mapping", func(t *testing.T) {
		got := (Resolver{}).ExpandPermissions([]string{"admin"}, nil)
		if len(got) != 1 || got[0] != "*" {
			t.Fatalf("unexpected fallback permissions: %v", got)
		}
	})
}
