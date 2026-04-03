package services

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"servify/apps/server/internal/config"
	"servify/apps/server/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	db     *gorm.DB
	config *config.Config
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
	Name     string
	Phone    string
	Role     string
}

type LoginInput struct {
	Username string
	Password string
}

type AuthResult struct {
	Token     string
	ExpiresIn int
	User      *models.User
	SessionID string
}

func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, config: cfg}
}

func (s *AuthService) Register(ctx context.Context, req RegisterInput) (*AuthResult, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, ErrInvalidAuthInput
	}

	var count int64
	s.db.WithContext(ctx).Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		return nil, ErrAuthUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "customer"
	}
	if role == "admin" {
		var total int64
		s.db.WithContext(ctx).Model(&models.User{}).Count(&total)
		if total > 0 {
			role = "customer"
		}
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hash),
		Name:     req.Name,
		Phone:    req.Phone,
		Role:     role,
		Status:   "active",
	}
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return nil, ErrAuthUserAlreadyExists
		}
		return nil, err
	}

	session, err := s.createAuthSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return s.buildAuthResult(ctx, user, session)
}

func (s *AuthService) Login(ctx context.Context, req LoginInput) (*AuthResult, error) {
	var user models.User
	if err := s.db.WithContext(ctx).Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuthInvalidCredentials
		}
		return nil, err
	}
	if user.Status != "active" {
		return nil, ErrAuthUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, ErrAuthInvalidCredentials
	}
	session, err := s.createAuthSession(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return s.buildAuthResult(ctx, &user, session)
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uint, sessionID string) (*AuthResult, error) {
	user, err := s.GetCurrentUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	session, err := s.rotateOrCreateAuthSession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	return s.buildAuthResult(ctx, user, session)
}

func (s *AuthService) buildAuthResult(ctx context.Context, user *models.User, session *models.UserAuthSession) (*AuthResult, error) {
	if user == nil {
		return nil, fmt.Errorf("user is required")
	}
	if session == nil {
		return nil, fmt.Errorf("auth session is required")
	}

	now := time.Now()
	token, err := createHS256JWT(map[string]interface{}{
		"iat":                   now.Unix(),
		"sub":                   user.ID,
		"user_id":               user.ID,
		"roles":                 []string{user.Role},
		"exp":                   now.Add(s.config.JWT.ExpiresIn).Unix(),
		"token_version":         user.TokenVersion,
		"session_id":            session.ID,
		"session_token_version": session.TokenVersion,
	}, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}
	_ = s.db.WithContext(ctx).Model(user).Update("last_login", now).Error
	return &AuthResult{
		Token:     token,
		ExpiresIn: int(s.config.JWT.ExpiresIn.Seconds()),
		User:      user,
		SessionID: session.ID,
	}, nil
}

func (s *AuthService) createAuthSession(ctx context.Context, userID uint) (*models.UserAuthSession, error) {
	if s.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	session := &models.UserAuthSession{
		ID:           newAuthSessionID(),
		UserID:       userID,
		Status:       "active",
		TokenVersion: 0,
	}
	if err := s.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AuthService) rotateOrCreateAuthSession(ctx context.Context, userID uint, sessionID string) (*models.UserAuthSession, error) {
	if strings.TrimSpace(sessionID) == "" {
		return s.createAuthSession(ctx, userID)
	}

	var session models.UserAuthSession
	if err := s.db.WithContext(ctx).First(&session, "id = ? AND user_id = ?", sessionID, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.createAuthSession(ctx, userID)
		}
		return nil, err
	}

	now := time.Now().UTC()
	updates := map[string]interface{}{
		"status":            "active",
		"token_version":     gorm.Expr("COALESCE(token_version, 0) + 1"),
		"last_refreshed_at": now,
		"revoked_at":        nil,
	}
	if err := s.db.WithContext(ctx).Model(&models.UserAuthSession{}).
		Where("id = ? AND user_id = ?", sessionID, userID).
		Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).First(&session, "id = ? AND user_id = ?", sessionID, userID).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

func newAuthSessionID() string {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return fmt.Sprintf("auth_%d", time.Now().UnixNano())
	}
	return "auth_" + hex.EncodeToString(buf[:])
}

func createHS256JWT(payload map[string]interface{}, secret string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	payloadJSON, _ := json.Marshal(payload)
	enc := func(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }

	h := enc(headerJSON)
	p := enc(payloadJSON)
	signing := h + "." + p

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(signing))
	sig := mac.Sum(nil)
	return signing + "." + enc(sig), nil
}
