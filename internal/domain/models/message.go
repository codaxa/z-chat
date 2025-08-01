// Package models contains the domain models for the chat application.
package models

import "time"

// Message represents a message in the chat application.
type Message struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Sender    string    `json:"sender" gorm:"type:varchar(255);not null;index" validate:"required"`
	Receiver  string    `json:"receiver" gorm:"type:varchar(255);index" validate:"nefield=Sender"`
	Content   string    `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;default:now();index"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;default:now()"`
}

// TableName returns the table name for the Message model
func (Message) TableName() string {
	return "messages"
}

// Validate checks the Message fields for validity.
func (m *Message) Validate() error {
	return validate.Struct(m)
}
