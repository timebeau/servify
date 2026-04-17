package auth

import "github.com/gin-gonic/gin"

const (
	ContextUserID        = "user_id"
	ContextUserIDRaw     = "user_id_raw"
	ContextTenantID      = "tenant_id"
	ContextWorkspaceID   = "workspace_id"
	ContextTokenType     = "token_type"
	ContextPrincipalType = "principal_type"
	ContextRoles         = "roles"
	ContextPermissions   = "permissions"
)

// Subject is the normalized auth context projected into gin.Context by AuthMiddleware.
type Subject struct {
	UserID        uint
	HasUserID     bool
	UserIDRaw     interface{}
	TenantID      string
	WorkspaceID   string
	TokenType     string
	PrincipalType string
	Roles         []string
	Permissions   []string
}

func IsInternalPrincipalType(principalType string) bool {
	switch principalType {
	case PrincipalService, PrincipalAdmin:
		return true
	default:
		return false
	}
}

func (s Subject) IsInternalPrincipal() bool {
	return IsInternalPrincipalType(s.PrincipalType)
}

// SubjectFromGin extracts the normalized auth subject from gin.Context.
func SubjectFromGin(c *gin.Context) Subject {
	if c == nil {
		return Subject{}
	}

	subject := Subject{
		TenantID:      c.GetString(ContextTenantID),
		WorkspaceID:   c.GetString(ContextWorkspaceID),
		TokenType:     c.GetString(ContextTokenType),
		PrincipalType: c.GetString(ContextPrincipalType),
		Roles:         getGrantedRoles(c),
		Permissions:   getGrantedPermissions(c),
	}

	if v, ok := c.Get(ContextUserID); ok {
		if userID, ok := v.(uint); ok {
			subject.UserID = userID
			subject.HasUserID = true
		}
	}
	if v, ok := c.Get(ContextUserIDRaw); ok {
		subject.UserIDRaw = v
	}
	return subject
}
