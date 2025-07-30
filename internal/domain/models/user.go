// Package models contains the domain models for the chat application.
package models

import "time"

// User represents a user in the chat application.
type User struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()" validate:"uuid4"`
	Username  string    `json:"username" gorm:"type:varchar(50);uniqueIndex;not null" validate:"required,min=3,max=50,alphanum"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null" validate:"required,email"`
	Password  string    `json:"password" gorm:"type:varchar(64);not null" validate:"required,sha256"`
	CreatedAt time.Time `json:"created_at" gorm:"default:now();autoCreateTime" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now();autoUpdateTime" validate:"required"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// Validate checks the User fields for validity.
func (u *User) Validate() error {
	return validate.Struct(u)
}
