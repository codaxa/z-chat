// Package models contains the domain models for the chat application.
package models

import (
	"context"
	"time"
)

// Message represents a message in the chat application.
type Message struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"uuid4"`
	Sender    string    `json:"sender" gorm:"type:varchar(255);not null;index" validate:"required"`
	Receiver  string    `json:"receiver" gorm:"type:varchar(255);index" validate:"nefield=Sender"`
	Content   string    `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at" gorm:"not null;autoCreateTime;index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null;autoUpdateTime"`
}

// TableName returns the table name for the Message model
func (Message) TableName() string {
	return "messages"
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
