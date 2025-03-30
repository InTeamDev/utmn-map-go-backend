package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/admin/entities"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/admin/repository/sqlc"
)

type AdminConverter interface {
	AdminSqlcToEntity(admin sqlc.GetAdminByUsernameRow) entities.Admin
	AdminsSqlcToEntity(admins []sqlc.GetAllAdminsRow) []entities.Admin
}

type AdminRepository struct {
	db             *sql.DB
	adminConverter AdminConverter
}

func NewAdminRepository(db *sql.DB, converter AdminConverter) *AdminRepository {
	return &AdminRepository{
		db:             db,
		adminConverter: converter,
	}
}

func (r *AdminRepository) Create(ctx context.Context, admin entities.AdminCreateRequest) (*entities.Admin, error) {
	q := sqlc.New(r.db)
	// TODO: add hash password
	createdAdmin, err := q.CreateAdmin(ctx, sqlc.CreateAdminParams{
		Username:     admin.Username,
		PasswordHash: admin.Password,
		Email:        sql.NullString{String: admin.Email, Valid: admin.Email != ""},
		FullName:     sql.NullString{String: admin.FullName, Valid: admin.FullName != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("create admin: %w", err)
	}

	result := entities.Admin{
		ID:       createdAdmin.ID,
		Username: createdAdmin.Username,
		Email:    createdAdmin.Email,
		FullName: createdAdmin.FullName,
	}
	return &result, nil
}

func (r *AdminRepository) GetByUsername(ctx context.Context, username string) (*entities.Admin, error) {
	q := sqlc.New(r.db)
	admin, err := q.GetAdminByUsername(ctx, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get admin by username: %w", err)
	}

	//TODO: add password hash
	entity := entities.Admin{
		ID:           admin.ID,
		Username:     admin.Username,
		PasswordHash: admin.PasswordHash,
		Email:        admin.Email,
		FullName:     admin.FullName,
		LastLogin:    admin.LastLogin,
	}
	return &entity, nil
}

func (r *AdminRepository) SetPasswordResetToken(ctx context.Context, email, token string, expiresAt time.Time) error {
	q := sqlc.New(r.db)
	err := q.SetPasswordResetToken(ctx, sqlc.SetPasswordResetTokenParams{
		Email:     sql.NullString{String: email, Valid: email != ""},
		Token:     sql.NullString{String: token, Valid: token != ""},
		ExpiresAt: sql.NullTime{Time: expiresAt, Valid: !expiresAt.IsZero()},
	})
	if err != nil {
		return fmt.Errorf("set password reset token: %w", err)
	}
	return nil
}

func (r *AdminRepository) GetByResetToken(ctx context.Context, token string) (*entities.Admin, error) {
	q := sqlc.New(r.db)
	admin, err := q.GetAdminByResetToken(ctx, sql.NullString{String: token, Valid: token != ""})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get admin by reset token: %w", err)
	}

	result := entities.Admin{
		ID:             admin.ID,
		Username:       admin.Username,
		Email:          admin.Email,
		TokenExpiresAt: admin.TokenExpiresAt,
	}
	return &result, nil
}

func (r *AdminRepository) UpdateStatus(ctx context.Context, id int32, isActive bool) error {
	q := sqlc.New(r.db)
	err := q.UpdateAdminStatus(ctx, sqlc.UpdateAdminStatusParams{
		ID:       id,
		IsActive: sql.NullBool{Bool: isActive, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("update admin status: %w", err)
	}
	return nil
}

func (r *AdminRepository) UpdateProfile(ctx context.Context, id int32, email, fullName string) error {
	q := sqlc.New(r.db)
	err := q.UpdateAdminProfile(ctx, sqlc.UpdateAdminProfileParams{
		ID:       id,
		Email:    sql.NullString{String: email, Valid: email != ""},
		FullName: sql.NullString{String: fullName, Valid: fullName != ""},
	})
	if err != nil {
		return fmt.Errorf("update admin profile: %w", err)
	}
	return nil
}
