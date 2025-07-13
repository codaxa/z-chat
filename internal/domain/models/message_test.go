// Package models contains the domain models for the chat application.
package models

import (
	"testing"
	"time"
)

type messageTest struct {
	ID         string
	Sender     string
	Receiver   string
	Content    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	IsPositive bool
}

func TestMessageValidation(t *testing.T) {
	createdAt, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")
	updatedAt, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z")

	tests := []messageTest{
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: true,
		},
		{
			ID:         "invalid-uuid",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},

		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "invalid-uuid",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "invalid-uuid",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Content:    "Hello, World!",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "",
			CreatedAt:  createdAt,
			UpdatedAt:  updatedAt,
			IsPositive: false,
		},
		{
			ID:         "6a387a08-e972-4fbf-9146-0a39510c6d5a",
			Sender:     "b514feb6-13bb-44d2-86f6-59d05bd338c6",
			Receiver:   "2b5fdfe1-1c9c-4c52-bfe8-60d4a056047e",
			Content:    "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse fringilla, ante vel faucibus egestas, est orci interdum nunc, at finibus orci nunc ac dolor. Etiam vitae commodo orci, id faucibus enim. Sed ac sagittis ex. In quis erat accumsan, vehicula augue vitae, pellentesque lacus. Vestibulum blandit iaculis nibh, sit amet auctor urna interdum a. Donec ante odio, euismod vehicula laoreet non, ornare at dui. Suspendisse molestie dignissim quam, vel porttitor turpis congue vitae. Aliquam et tortor bibendum enim ultricies faucibus at nec diam. Maecenas lobortis lacinia sem sed dapibus. Praesent gravida vel nulla a placerat. Vivamus ultrices tempus ultricies. Phasellus eu sem non diam commodo pellentesque. Sed ultrices pretium iaculis. Sed accumsan mi sed diam tincidunt pellentesque. Praesent sollicitudin in lectus at volutpat. Phasellus efficitur velit et interdum cursus. Integer ac elit ipsum. Praesent ac dignissim ipsum. Aenean quis efficitur ipsum. Nunc maximus aliquet nibh eget vel.",
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
				Receiver:  test.Receiver,
				CreatedAt: test.CreatedAt,
				UpdatedAt: test.UpdatedAt,
			}

			err := message.Validate()
			if test.IsPositive && err != nil {
				t.Errorf("expected no error for:\n\t\t%+v \nbut, got the following error:\n\t\t%v", test, err)
			} else if !test.IsPositive && err == nil {
				t.Error("expected error for invalid message, got nil")
			}
		})
	}
}
