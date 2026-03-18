package auth

import (
	"net/http"
	"strings"
	"time"

	"servify/apps/server/internal/config"

	"github.com/gin-gonic/gin"
)

type MiddlewareConfig struct {
	Secret string
	RBAC   config.RBACConfig
	Now    func() time.Time
}

func MiddlewareConfigFromApp(cfg *config.Config) MiddlewareConfig {
	var out MiddlewareConfig
	if cfg == nil {
		return out
	}
	out.Secret = cfg.JWT.Secret
	out.RBAC = cfg.Security.RBAC
	return out
}

// AuthMiddleware enforces Bearer JWT auth and projects normalized claims into gin.Context.
func AuthMiddleware(cfg MiddlewareConfig) gin.HandlerFunc {
	validator := Validator{
		Secret: cfg.Secret,
		Now:    cfg.Now,
	}
	resolver := Resolver{RBAC: cfg.RBAC}

	return func(c *gin.Context) {
		ah := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(ah), "bearer ") {
			abortJSON(c, http.StatusUnauthorized, "Unauthorized", "missing bearer token")
			return
		}
		token := strings.TrimSpace(ah[len("Bearer "):])
		if token == "" || validator.Secret == "" {
			abortJSON(c, http.StatusUnauthorized, "Unauthorized", "invalid token or server misconfig")
			return
		}

		payload, err := validator.ValidateToken(token)
		if err != nil {
			abortJSON(c, http.StatusUnauthorized, "Unauthorized", err.Error())
			return
		}

		claims := extractClaims(payload, resolver)
		if claims.HasUserID {
			c.Set("user_id", claims.UserID)
		}
		if claims.UserIDRaw != nil {
			c.Set("user_id_raw", claims.UserIDRaw)
		}
		if len(claims.Roles) > 0 {
			c.Set("roles", claims.Roles)
		}
		if len(claims.Permissions) > 0 {
			c.Set("permissions", claims.Permissions)
		}

		c.Next()
	}
}

// RequirePermissionsAny requires the caller to have at least one listed permission.
func RequirePermissionsAny(required ...string) gin.HandlerFunc {
	req := make([]string, 0, len(required))
	for _, r := range required {
		if s := strings.TrimSpace(r); s != "" {
			req = append(req, s)
		}
	}
	return func(c *gin.Context) {
		granted := getGrantedPermissions(c)
		for _, r := range req {
			if HasPermission(granted, r) {
				c.Next()
				return
			}
		}
		abortJSON(c, http.StatusForbidden, "Forbidden", "insufficient permission")
	}
}

// RequirePermissionsAll requires the caller to have all listed permissions.
func RequirePermissionsAll(required ...string) gin.HandlerFunc {
	req := make([]string, 0, len(required))
	for _, r := range required {
		if s := strings.TrimSpace(r); s != "" {
			req = append(req, s)
		}
	}
	return func(c *gin.Context) {
		granted := getGrantedPermissions(c)
		for _, r := range req {
			if !HasPermission(granted, r) {
				abortJSON(c, http.StatusForbidden, "Forbidden", "insufficient permission")
				return
			}
		}
		c.Next()
	}
}

// RequireResourcePermission maps HTTP methods to resource permissions.
func RequireResourcePermission(resource string) gin.HandlerFunc {
	resource = strings.TrimSpace(resource)
	return func(c *gin.Context) {
		perm := ResourcePermission(resource, c.Request.Method)
		RequirePermissionsAny(perm, resource+".*", "*")(c)
	}
}

// RequireRolesAny requires at least one matching role from the normalized context.
func RequireRolesAny(required ...string) gin.HandlerFunc {
	reqSet := make(map[string]struct{}, len(required))
	for _, r := range required {
		if role := strings.TrimSpace(r); role != "" {
			reqSet[role] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		for _, r := range getGrantedRoles(c) {
			if _, ok := reqSet[r]; ok {
				c.Next()
				return
			}
		}
		abortJSON(c, http.StatusForbidden, "Forbidden", "insufficient role")
	}
}

func getGrantedPermissions(c *gin.Context) []string {
	if v, ok := c.Get("permissions"); ok {
		if perms, ok := v.([]string); ok {
			return perms
		}
	}
	return nil
}

func getGrantedRoles(c *gin.Context) []string {
	if v, ok := c.Get("roles"); ok {
		switch t := v.(type) {
		case []string:
			return t
		case []interface{}:
			var out []string
			for _, it := range t {
				if s, ok := it.(string); ok {
					out = append(out, s)
				}
			}
			return out
		case string:
			if t != "" {
				return []string{t}
			}
		}
	}
	return nil
}

func abortJSON(c *gin.Context, status int, errCode, msg string) {
	c.AbortWithStatusJSON(status, gin.H{
		"error":   errCode,
		"message": msg,
	})
}
