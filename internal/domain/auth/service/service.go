package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/auth/entity"
	userEntity "github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/entity"
	userService "github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Errors returned by the auth service
var (
	ErrHashNotFound = errors.New("registration hash not found")
	ErrHashExpired  = errors.New("registration hash expired")
	ErrCodeNotFound = errors.New("authentication code not found")
	ErrCodeExpired  = errors.New("authentication code expired")
	ErrInvalidCode  = errors.New("invalid authentication code")
	ErrUserNotFound = errors.New("user not found")
)

// TokenClaims represents JWT claims
type TokenClaims struct {
	UserID     uuid.UUID       `json:"user_id"`
	TelegramID int64           `json:"telegram_id"`
	Username   string          `json:"username"`
	Role       userEntity.Role `json:"role"`
	jwt.RegisteredClaims
}

// AuthService provides authentication-related operations
type AuthService struct {
	userService        *userService.Service
	regHashStore       map[string]*entity.RegistrationHash
	authCodeStore      map[string]*entity.AuthCode
	mu                 sync.RWMutex
	jwtSecret          []byte
	jwtExpirationHours int
}

// NewAuthService creates a new authentication service
func NewAuthService(userService *userService.Service, jwtSecret string, jwtExpirationHours int) *AuthService {
	service := &AuthService{
		userService:        userService,
		regHashStore:       make(map[string]*entity.RegistrationHash),
		authCodeStore:      make(map[string]*entity.AuthCode),
		jwtSecret:          []byte(jwtSecret),
		jwtExpirationHours: jwtExpirationHours,
	}

	// Run cleanup goroutine
	go service.cleanupExpiredItems()

	return service
}

// GenerateRegistrationHash creates a new registration hash
func (s *AuthService) GenerateRegistrationHash() *entity.RegistrationHash {
	s.mu.Lock()
	defer s.mu.Unlock()

	regHash := entity.NewRegistrationHash()
	s.regHashStore[regHash.Hash] = regHash
	return regHash
}

// ValidateRegistrationHash validates a registration hash
func (s *AuthService) ValidateRegistrationHash(hash string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	regHash, exists := s.regHashStore[hash]
	if !exists {
		return ErrHashNotFound
	}

	if regHash.IsExpired() {
		return ErrHashExpired
	}

	return nil
}

// ConsumeRegistrationHash validates and then removes a registration hash
func (s *AuthService) ConsumeRegistrationHash(hash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	regHash, exists := s.regHashStore[hash]
	if !exists {
		return ErrHashNotFound
	}

	if regHash.IsExpired() {
		delete(s.regHashStore, hash)
		return ErrHashExpired
	}

	delete(s.regHashStore, hash)
	return nil
}

// GenerateAuthCode creates a new authentication code for a user
func (s *AuthService) GenerateAuthCode(ctx context.Context, username string) (*entity.AuthCode, error) {
	// Verify user exists
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	authCode := entity.NewAuthCode(user.Username)
	s.authCodeStore[authCode.Username] = authCode
	return authCode, nil
}

// ValidateAuthCode validates an authentication code
func (s *AuthService) ValidateAuthCode(ctx context.Context, username, code string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	authCode, exists := s.authCodeStore[username]
	if !exists {
		return ErrCodeNotFound
	}

	if authCode.IsExpired() {
		return ErrCodeExpired
	}

	if authCode.Code != code {
		return ErrInvalidCode
	}

	return nil
}

// ConsumeAuthCode validates and then removes an authentication code
func (s *AuthService) ConsumeAuthCode(ctx context.Context, username, code string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	authCode, exists := s.authCodeStore[username]
	if !exists {
		return ErrCodeNotFound
	}

	if authCode.IsExpired() {
		delete(s.authCodeStore, username)
		return ErrCodeExpired
	}

	if authCode.Code != code {
		return ErrInvalidCode
	}

	delete(s.authCodeStore, username)
	return nil
}

// GenerateJWT generates a JWT token for a user
func (s *AuthService) GenerateJWT(ctx context.Context, username string) (string, error) {
	// Get user
	user, err := s.userService.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, userService.ErrUserNotFound) {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Create token with claims
	expiresAt := time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour)
	claims := TokenClaims{
		UserID:     user.ID,
		TelegramID: user.TelegramID,
		Username:   user.Username,
		Role:       user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// VerifyJWT verifies a JWT token and returns the claims
func (s *AuthService) VerifyJWT(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token claims")
}

// cleanupExpiredItems periodically removes expired registration hashes and auth codes
func (s *AuthService) cleanupExpiredItems() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()

		// Clean up expired registration hashes
		for hash, regHash := range s.regHashStore {
			if regHash.IsExpired() {
				delete(s.regHashStore, hash)
			}
		}

		// Clean up expired auth codes
		for username, authCode := range s.authCodeStore {
			if authCode.IsExpired() {
				delete(s.authCodeStore, username)
			}
		}

		s.mu.Unlock()
	}
}
