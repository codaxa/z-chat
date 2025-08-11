// Package models contains the domain models for the chat application.
package models

import (
	"testing"
	"time"
)

type messageTest struct {
	ID         string
	Sender     string
	Content    string
	RoomID     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsPositive bool
}

func TestMessageValidation(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	testRoom := Room{
		ID:        "room-123",
		Name:      "Test Room",
		CreatedBy: "test-user-id",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	tests := []messageTest{
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Content:    "Hello, World!",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "invalid-uuid",
			Content:    "Hello, World!",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "",
			RoomID:     "room-123",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			RoomID:     "",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
	}

	for _, test := range tests {
		t.Run(test.ID, func(t *testing.T) {
			message := Message{
				ID:        test.ID,
				Content:   test.Content,
				Sender:    test.Sender,
				RoomID:    test.RoomID,
				CreatedAt: test.CreatedAt,
				UpdatedAt: test.UpdatedAt,
				Room:      testRoom,
			}

			err := message.Validate()
			if test.IsPositive && err != nil {
				t.Errorf("expected no error for:\n\t\t%+v \nbut, got the following error:\n\t\t%v", test, err)
			} else if !test.IsPositive && err == nil {
				t.Errorf("expected error for:\n\t\t%+v \nbut, got no error", test)
			}
		})
	}
}
