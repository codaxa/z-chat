// Package models contains the domain models for the chat application.
package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Create a single instance of the validator
var validate = validator.New()

// User represents a user in the chat application.
type User struct {
	ID             string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"uuid4"`
	Username       string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=3,max=50,alphanum"`
	Email          string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	HashedPassword string    `json:"password" gorm:"type:varchar(64);not null" validate:"required,sha256"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;autoCreateTime" validate:"required,gt=0001-01-01"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null;autoUpdateTime" validate:"required,gt=0001-01-01"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// Validate checks the User fields for validity.
func (u *User) Validate() error {
	return validate.Struct(u)
}
