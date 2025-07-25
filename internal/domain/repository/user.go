// Package repository defines interfaces for data access operations
package repository

import (
	"context"
	"z-chat/internal/domain/models"
)

// UserRepo defines operations for managing users in the data store
type UserRepo interface {
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Add(ctx context.Context, u models.User) error
}
