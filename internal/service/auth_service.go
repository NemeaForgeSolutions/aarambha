package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aarambha-agency/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionCookieName = "aarambha_session"
	DefaultSessionTTL = 7 * 24 * time.Hour
)

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
	Phone    string
	Role     string
}

type AuthService struct {
	userRepo      domain.UserRepository
	sessionSecret []byte
}

func NewAuthService(repo domain.UserRepository, sessionSecret string) *AuthService {
	if strings.TrimSpace(sessionSecret) == "" {
		sessionSecret = "aarambha-default-secret-salt-2026"
	}
	return &AuthService{
		userRepo:      repo,
		sessionSecret: []byte(sessionSecret),
	}
}

func (s *AuthService) Register(req RegisterRequest) (*domain.User, error) {
	name := strings.TrimSpace(req.Name)
	email := strings.ToLower(strings.TrimSpace(req.Email))
	password := req.Password

	if name == "" {
		return nil, errors.New("full name is required")
	}
	if email == "" || !strings.Contains(email, "@") {
		return nil, errors.New("valid email address is required")
	}
	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	role := strings.ToLower(strings.TrimSpace(req.Role))
	if role != domain.RoleAgent {
		role = domain.RoleUser
	}

	// Check if already registered
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	newUser := &domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
		Phone:        strings.TrimSpace(req.Phone),
		IsActive:     true,
	}

	if err := s.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return newUser.Sanitized(), nil
}

func (s *AuthService) Authenticate(email, password string) (*domain.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || password == "" {
		return nil, domain.ErrInvalidCredentials
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Login is restricted to admin and agents for now
	if !user.IsAdmin() && !user.IsAgent() {
		return nil, domain.ErrLoginRestricted
	}

	_ = s.userRepo.UpdateLastLogin(user.ID)

	return user.Sanitized(), nil
}

func (s *AuthService) CreateSessionToken(userID string, ttl time.Duration) string {
	if ttl == 0 {
		ttl = DefaultSessionTTL
	}
	expiresAt := time.Now().Add(ttl).Unix()
	payload := fmt.Sprintf("%s:%d", userID, expiresAt)

	mac := hmac.New(sha256.New, s.sessionSecret)
	mac.Write([]byte(payload))
	sig := hex.EncodeToString(mac.Sum(nil))

	raw := fmt.Sprintf("%s.%s", payload, sig)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func (s *AuthService) ValidateSessionToken(tokenStr string) (*domain.User, error) {
	if strings.TrimSpace(tokenStr) == "" {
		return nil, errors.New("empty session token")
	}

	raw, err := base64.URLEncoding.DecodeString(tokenStr)
	if err != nil {
		return nil, errors.New("malformed token encoding")
	}

	parts := strings.Split(string(raw), ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token structure")
	}

	payload := parts[0]
	expectedSig := parts[1]

	mac := hmac.New(sha256.New, s.sessionSecret)
	mac.Write([]byte(payload))
	actualSig := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expectedSig), []byte(actualSig)) {
		return nil, errors.New("invalid token signature")
	}

	payloadParts := strings.Split(payload, ":")
	if len(payloadParts) != 2 {
		return nil, errors.New("invalid token payload")
	}

	userID := payloadParts[0]
	expiresAtUnix, err := strconv.ParseInt(payloadParts[1], 10, 64)
	if err != nil {
		return nil, errors.New("invalid expiration in token")
	}

	if time.Now().Unix() > expiresAtUnix {
		return nil, errors.New("session has expired")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, domain.ErrUserInactive
	}

	return user.Sanitized(), nil
}

func (s *AuthService) SetSessionCookie(w http.ResponseWriter, user *domain.User, isProduction bool) {
	token := s.CreateSessionToken(user.ID, DefaultSessionTTL)

	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(DefaultSessionTTL),
		MaxAge:   int(DefaultSessionTTL.Seconds()),
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *AuthService) ClearSessionCookie(w http.ResponseWriter, isProduction bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isProduction,
		SameSite: http.SameSiteLaxMode,
	})
}

func (s *AuthService) GetUserFromRequest(r *http.Request) (*domain.User, error) {
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil {
		return nil, err
	}
	return s.ValidateSessionToken(cookie.Value)
}

func (s *AuthService) GetUserByID(id string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return user.Sanitized(), nil
}

func (s *AuthService) UpdateProfile(userID, name, phone string) (*domain.User, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(name) != "" {
		user.Name = strings.TrimSpace(name)
	}
	user.Phone = strings.TrimSpace(phone)

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user.Sanitized(), nil
}
