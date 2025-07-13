// Package models contains the domain models for the chat application.
package models

import "time"

// User represents a user in the chat application.
type User struct {
	ID        string    `json:"id" validate:"uuid4"`
	Username  string    `json:"username" validate:"required,min=3,max=50"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" validate:"required,sha256"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate checks the User fields for validity.
func (u *User) Validate() error {
	return validate.Struct(u)
}
