// Package models contains the domain models for the chat application.
package models

import (
	"strings"
	"testing"
	"time"
)

// TestUserValidation tests various user validation scenarios
func TestUserValidation(t *testing.T) {
	createdAt, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse createdAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse updatedAt: %v", err)
	}

	validID := "6a387a08-e972-4fbf-9146-0a39510c6d5a"
	validUsername := "testuser"
	validEmail := "test.user@email.com"
	validPassword := "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"

	tests := []struct {
		name        string
		user        User
		expectError bool
		errorField  string
	}{
		{
			name: "valid user",
			user: User{
				ID:        validID,
				Username:  validUsername,
				Email:     validEmail,
				Password:  validPassword,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectError: false,
		},
		{
			name: "empty ID",
			user: User{
				ID:        "",
				Username:  validUsername,
				Email:     validEmail,
				Password:  validPassword,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectError: true,
			errorField:  "ID",
		},
		{
			name: "special characters in username",
			user: User{
				ID:        validID,
				Username:  "test@user",
				Email:     validEmail,
				Password:  validPassword,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			expectError: true,
			errorField:  "Username",
		},
		{
			name: "zero time for CreatedAt",
			user: User{
				ID:        validID,
				Username:  validUsername,
				Email:     validEmail,
				Password:  validPassword,
				CreatedAt: time.Time{},
				UpdatedAt: updatedAt,
			},
			expectError: true,
			errorField:  "CreatedAt",
		},
		{
			name: "zero time for UpdatedAt",
			user: User{
				ID:        validID,
				Username:  validUsername,
				Email:     validEmail,
				Password:  validPassword,
				CreatedAt: createdAt,
				UpdatedAt: time.Time{},
			},
			expectError: true,
			errorField:  "UpdatedAt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.expectError && err == nil {
				t.Errorf("expected error for %s, got nil", tt.name)
			}

			if !tt.expectError && err != nil {
				t.Errorf("expected no error for %s, but got: %v", tt.name, err)
			}

			if tt.expectError && err != nil && tt.errorField != "" {
				if !containsField(err.Error(), tt.errorField) {
					t.Errorf("expected error for field %s, but got: %v", tt.errorField, err)
				}
			}
		})
	}
}

// Helper function to check if error message contains field name
func containsField(errorMsg, fieldName string) bool {
	return strings.Contains(errorMsg, fieldName)
}
