// Package models contains the domain models for the chat application.
package models

import "time"

// Message represents a message in the chat application.
type Message struct {
	ID        string    `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content" validate:"required,min=1,max=1000"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks the Message fields for validity.
func (m *Message) Validate() error {
	return validate.Struct(m)
}
