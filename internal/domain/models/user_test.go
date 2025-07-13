// Package models contains the domain models for the chat application.
package models

import (
	"testing"
	"time"
)

type userTest struct {
	ID         string
	Username   string
	Email      string
	Password   string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsPositive bool
}

func TestUserValidation(t *testing.T) {
	createdAt, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse createdAt: %v", err)
	}
	updatedAt, err := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse updatedAt: %v", err)
	}

	tests := []userTest{
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "testuser",
			Email:      "test.user@email.com",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "invalid-uuid",
			Username:   "testuser",
			Email:      "test.user@email.com",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Email:      "test.user@email.com",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "tu",
			Email:      "test.user@email.com",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "test_user_with_an_exceptionally_long_username_example_123",
			Email:      "test.user@email.com",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "testuser",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "testuser",
			Email:      "test-user-email",
			Password:   "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "testuser",
			Email:      "test.user@email.com",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Username:   "testuser",
			Email:      "test.user@email.com",
			Password:   "unhashesdpassword",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
	}

	for _, test := range tests {
		t.Run(test.ID, func(t *testing.T) {
			user := User{
				ID:        test.ID,
				Username:  test.Username,
				Email:     test.Email,
				Password:  test.Password,
				CreatedAt: test.CreatedAt,
				UpdatedAt: test.UpdatedAt,
			}

			err := user.Validate()
			if test.IsPositive && err != nil {
				t.Errorf("expected no error for:\n\t\t%+v \nbut, got the following error:\n\t\t%v", test, err)
			} else if !test.IsPositive && err == nil {
				t.Error("expected error for invalid user, got nil")
			}
		})
	}
}
