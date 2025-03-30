package entities

import (
	"database/sql"
)

type Admin struct {
	ID                 int32          `json:"id"`
	Username           string         `json:"username"`
	PasswordHash       string         `json:"-"`
	Email              sql.NullString `json:"email"`
	FullName           sql.NullString `json:"full_name"`
	LastLogin          sql.NullTime   `json:"last_login,omitempty"`
	IsActive           bool           `json:"is_active"`
	Salt               string         `json:"-"`
	PasswordResetToken sql.NullString `json:"-"`
	TokenExpiresAt     sql.NullTime   `json:"-"`
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
