package services

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
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

	return s.buildAuthResult(ctx, user)
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
	return s.buildAuthResult(ctx, &user)
}

func (s *AuthService) GetCurrentUser(ctx context.Context, userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, userID uint) (*AuthResult, error) {
	user, err := s.GetCurrentUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.buildAuthResult(ctx, user)
}

func (s *AuthService) buildAuthResult(ctx context.Context, user *models.User) (*AuthResult, error) {
	token, err := createHS256JWT(map[string]interface{}{
		"iat":     time.Now().Unix(),
		"sub":     user.ID,
		"user_id": user.ID,
		"roles":   []string{user.Role},
		"exp":     time.Now().Add(s.config.JWT.ExpiresIn).Unix(),
	}, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_ = s.db.WithContext(ctx).Model(user).Update("last_login", now).Error
	return &AuthResult{
		Token:     token,
		ExpiresIn: int(s.config.JWT.ExpiresIn.Seconds()),
		User:      user,
	}, nil
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
