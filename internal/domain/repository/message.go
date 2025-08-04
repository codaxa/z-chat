// Package repository defines interfaces for data access operations
package repository

import (
	"context"
	"z-chat/internal/domain/models"
)

// MessageRepository defines the interface for message-related operations.
type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *models.Message) error
	GetMessageByID(ctx context.Context, id string) (*models.Message, error)
	GetMessagesByRoom(ctx context.Context, roomID string) ([]*models.Message, error)
}
