package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/entity"
	"github.com/google/uuid"
)

// Errors returned by the repository
var (
	ErrUserNotFound = errors.New("user not found")
)

// PostgresRepository implements Repository interface for PostgreSQL
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create adds a new user to the repository
func (r *PostgresRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (id, telegram_id, username, role, created_at, modified_at, photo_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.TelegramID,
		user.Username,
		user.Role,
		user.CreatedAt,
		user.ModifiedAt,
		user.PhotoURL,
	)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by ID
func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `
		SELECT id, telegram_id, username, role, created_at, modified_at, photo_url
		FROM users
		WHERE id = $1
	`
	return r.queryUser(ctx, query, id)
}

// GetByTelegramID retrieves a user by Telegram ID
func (r *PostgresRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	query := `
		SELECT id, telegram_id, username, role, created_at, modified_at, photo_url
		FROM users
		WHERE telegram_id = $1
	`
	return r.queryUser(ctx, query, telegramID)
}

// GetByUsername retrieves a user by username
func (r *PostgresRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, telegram_id, username, role, created_at, modified_at, photo_url
		FROM users
		WHERE username = $1
	`
	return r.queryUser(ctx, query, username)
}

// UpdateRole updates the role of a user
func (r *PostgresRepository) UpdateRole(ctx context.Context, id uuid.UUID, role entity.Role) error {
	query := `
		UPDATE users
		SET role = $1, modified_at = $2
		WHERE id = $3
	`
	result, err := r.db.ExecContext(ctx, query, role, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Update updates user information
func (r *PostgresRepository) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET username = $1, photo_url = $2, modified_at = $3
		WHERE id = $4
	`
	user.ModifiedAt = time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.PhotoURL,
		user.ModifiedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// List retrieves all users
func (r *PostgresRepository) List(ctx context.Context) ([]*entity.User, error) {
	query := `
		SELECT id, telegram_id, username, role, created_at, modified_at, photo_url
		FROM users
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		var photoURL sql.NullString
		err := rows.Scan(
			&user.ID,
			&user.TelegramID,
			&user.Username,
			&user.Role,
			&user.CreatedAt,
			&user.ModifiedAt,
			&photoURL,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		if photoURL.Valid {
			user.PhotoURL = photoURL.String
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// queryUser executes a query that returns a single user
func (r *PostgresRepository) queryUser(ctx context.Context, query string, args ...interface{}) (*entity.User, error) {
	var user entity.User
	var photoURL sql.NullString
	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.TelegramID,
		&user.Username,
		&user.Role,
		&user.CreatedAt,
		&user.ModifiedAt,
		&photoURL,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to query user: %w", err)
	}
	if photoURL.Valid {
		user.PhotoURL = photoURL.String
	}
	return &user, nil
} 