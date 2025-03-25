package entities

import (
	"database/sql"
	"time"
)

type Admin struct {
	ID        int32          `json:"id"`
	Username  string         `json:"username"`
	Email     sql.NullString `json:"email"`
	FullName  sql.NullString `json:"full_name"`
	CreatedAt time.Time      `json:"created_at"`
	LastLogin sql.NullTime   `json:"last_login,omitempty"`
	TokenExp  sql.NullTime   `json:"-"`
}

type AdminWithCredentials struct {
	Admin
	PasswordHash string `json:"-"`
	Salt         string `json:"-"`
}

type AdminCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

type AdminUpdateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"full_name" validate:"required"`
}

type AdminPasswordUpdateRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type AdminPasswordResetRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type AdminLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Admin
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
