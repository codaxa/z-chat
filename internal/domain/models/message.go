// Package models contains the domain models for the chat application.
package models

import (
	"context"
	"time"
)

// Message represents a message in the chat application.
type Message struct {
	ID        string    `json:"id" validate:"uuid4" gorm:"primaryKey"`
	Sender    string    `json:"sender" validate:"required"`
	Receiver  string    `json:"receiver" validate:"nefield=Sender"`
	Content   string    `json:"content" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// MessageRepository defines the interface for message-related operations.
type MessageRepository interface {
	CreateMessage(ctx context.Context, msg *Message) error
	GetMessageByID(ctx context.Context, id string) (*Message, error)
}

// Validate checks the Message fields for validity.
func (m *Message) Validate() error {
	return validate.Struct(m)
}
