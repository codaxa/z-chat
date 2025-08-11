// Package repository defines interfaces for data access operations
package repository

import (
	"context"
	"z-chat/internal/domain/models"
)

// UserRepository defines operations for managing users in the data store
type UserRepository interface {
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, u models.User) error
}
