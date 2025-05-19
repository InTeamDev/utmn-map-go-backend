package entity

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user roles in the system
type Role string

const (
	// RoleUser is the base user role
	RoleUser Role = "user"
	// RoleAdmin is the admin role with access to Admin Panel
	RoleAdmin Role = "admin"
	// RoleCurator is the role that can grant/revoke Admin rights
	RoleCurator Role = "curator"
)

// User represents a user in the system
type User struct {
	ID          uuid.UUID `json:"id"`
	TelegramID  int64     `json:"telegram_id"`
	Username    string    `json:"username"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	ModifiedAt  time.Time `json:"modified_at"`
	PhotoURL    string    `json:"photo_url,omitempty"`
}

// NewUser creates a new User with RoleUser
func NewUser(telegramID int64, username string, photoURL string) *User {
	now := time.Now()
	return &User{
		ID:          uuid.New(),
		TelegramID:  telegramID,
		Username:    username,
		Role:        RoleUser,
		CreatedAt:   now,
		ModifiedAt:  now,
		PhotoURL:    photoURL,
	}
} 