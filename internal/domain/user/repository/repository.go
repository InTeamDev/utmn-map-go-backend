package repository

import (
	"context"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/entity"
	"github.com/google/uuid"
)

// Repository defines the interface for user data operations
type Repository interface {
	// Create adds a new user to the repository
	Create(ctx context.Context, user *entity.User) error
	
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	
	// GetByTelegramID retrieves a user by Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error)
	
	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	
	// UpdateRole updates the role of a user
	UpdateRole(ctx context.Context, id uuid.UUID, role entity.Role) error
	
	// Update updates user information
	Update(ctx context.Context, user *entity.User) error
	
	// List retrieves all users
	List(ctx context.Context) ([]*entity.User, error)
} 