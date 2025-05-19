package entity

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	// RegistrationHashTTL is the time-to-live for registration hashes (5 minutes)
	RegistrationHashTTL = 5 * time.Minute

	// AuthCodeTTL is the time-to-live for authentication codes (5 minutes)
	AuthCodeTTL = 5 * time.Minute

	// AuthCodeLength is the length of the authentication code
	AuthCodeLength = 6
)

// RegistrationHash represents a registration hash for Telegram deep-link
type RegistrationHash struct {
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired checks if the registration hash has expired
func (r *RegistrationHash) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

// NewRegistrationHash creates a new registration hash
func NewRegistrationHash() *RegistrationHash {
	now := time.Now()
	return &RegistrationHash{
		Hash:      uuid.New().String(),
		CreatedAt: now,
		ExpiresAt: now.Add(RegistrationHashTTL),
	}
}

// AuthCode represents a one-time code for authentication
type AuthCode struct {
	Username  string    `json:"username"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired checks if the auth code has expired
func (c *AuthCode) IsExpired() bool {
	return time.Now().After(c.ExpiresAt)
}

// NewAuthCode creates a new authentication code for a user
func NewAuthCode(username string) *AuthCode {
	now := time.Now()
	return &AuthCode{
		Username:  username,
		Code:      generateNumericCode(AuthCodeLength),
		CreatedAt: now,
		ExpiresAt: now.Add(AuthCodeTTL),
	}
}

// generateNumericCode generates a random numeric code of the specified length
func generateNumericCode(length int) string {
	const digits = "0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = digits[rand.Intn(len(digits))]
	}
	return string(b)
}
