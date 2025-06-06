package entities

import "time"

type User struct {
	ID        string
	TGID      int64
	Username  string
	Roles     []string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AuthCode struct {
	Code      string
	ExpiresAt time.Time
	Attempts  int
	SentAt    time.Time
}

type RefreshToken struct {
	JTI       string
	UserID    string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Revoked   bool
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
