package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/entity"
	"github.com/InTeamDev/utmn-map-go-backend/internal/domain/user/repository"
	"github.com/google/uuid"
)

// Errors returned by the service
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInsufficientRights = errors.New("insufficient rights")
)

// Service provides user-related operations
type Service struct {
	repo repository.Repository
}

// NewService creates a new user service
func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// RegisterUser registers a new user
func (s *Service) RegisterUser(ctx context.Context, telegramID int64, username, photoURL string) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		// User already exists
		return existingUser, nil
	}
	if !errors.Is(err, repository.ErrUserNotFound) {
		return nil, fmt.Errorf("failed to check if user exists: %w", err)
	}

	// Create new user
	user := entity.NewUser(telegramID, username, photoURL)
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// GetUserByTelegramID retrieves a user by Telegram ID
func (s *Service) GetUserByTelegramID(ctx context.Context, telegramID int64) (*entity.User, error) {
	user, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by telegram ID: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (s *Service) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

// PromoteUser promotes a user to admin role
func (s *Service) PromoteUser(ctx context.Context, curatorID uuid.UUID, username string) (*entity.User, error) {
	// Verify curator has sufficient rights
	curator, err := s.GetUser(ctx, curatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get curator: %w", err)
	}
	if curator.Role != entity.RoleCurator {
		return nil, ErrInsufficientRights
	}

	// Get user to promote
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Can only promote User to Admin
	if user.Role != entity.RoleUser {
		return nil, fmt.Errorf("can only promote users with User role")
	}

	// Update role
	if err := s.repo.UpdateRole(ctx, user.ID, entity.RoleAdmin); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	// Retrieve updated user
	return s.GetUser(ctx, user.ID)
}

// DemoteUser demotes a user from admin role
func (s *Service) DemoteUser(ctx context.Context, curatorID uuid.UUID, username string) (*entity.User, error) {
	// Verify curator has sufficient rights
	curator, err := s.GetUser(ctx, curatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get curator: %w", err)
	}
	if curator.Role != entity.RoleCurator {
		return nil, ErrInsufficientRights
	}

	// Get user to demote
	user, err := s.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// Can only demote Admin to User
	if user.Role != entity.RoleAdmin {
		return nil, fmt.Errorf("can only demote users with Admin role")
	}

	// Update role
	if err := s.repo.UpdateRole(ctx, user.ID, entity.RoleUser); err != nil {
		return nil, fmt.Errorf("failed to update user role: %w", err)
	}

	// Retrieve updated user
	return s.GetUser(ctx, user.ID)
}

// ListUsers retrieves all users
func (s *Service) ListUsers(ctx context.Context) ([]*entity.User, error) {
	return s.repo.List(ctx)
}
