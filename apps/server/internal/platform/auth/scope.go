package auth

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	ScopeModeGlobal    = "global"
	ScopeModeTenant    = "tenant"
	ScopeModeWorkspace = "workspace"
)

// Scope is the normalized tenant/workspace view derived from the auth subject.
type Scope struct {
	TenantID    string
	WorkspaceID string
}

func ScopeFromSubject(subject Subject) Scope {
	return Scope{
		TenantID:    subject.TenantID,
		WorkspaceID: subject.WorkspaceID,
	}
}

func ScopeFromGin(c *gin.Context) Scope {
	return ScopeFromSubject(SubjectFromGin(c))
}

func (s Scope) HasTenant() bool {
	return s.TenantID != ""
}

func (s Scope) HasWorkspace() bool {
	return s.WorkspaceID != ""
}

func (s Scope) Mode() string {
	switch {
	case s.HasWorkspace():
		return ScopeModeWorkspace
	case s.HasTenant():
		return ScopeModeTenant
	default:
		return ScopeModeGlobal
	}
}

func (s Scope) Validate() error {
	if s.HasWorkspace() && !s.HasTenant() {
		return fmt.Errorf("workspace scope requires tenant scope")
	}
	return nil
}
