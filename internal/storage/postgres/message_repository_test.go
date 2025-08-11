package postgres

import (
	"context"
	"errors"
	"testing"
	"time"
	"z-chat/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

// Test for NewMessageRepository function
func TestNewMessageRepository(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := NewMessageRepository(mock)

	if repo == nil {
		t.Fatal("expected MessageRepository instance, got nil")
	}
	if repo.db != mock {
		t.Error("expected repository to have correct database connection reference")
	}
}

// Test for CreateMessage method
func TestMessageRepository_CreateMessage(t *testing.T) {
	tests := []struct {
		name    string
		message *models.Message
		mockErr error
		wantErr bool
	}{
		{
			name: "successful creation",
			message: &models.Message{
				ID:        "550e8400-e29b-41d4-a716-446655440000",
				Sender:    "user1",
				Content:   "Hello, World!",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				RoomID:    "some-room-id",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "database error",
			message: &models.Message{
				ID:        "550e8400-e29b-41d4-a716-446655440000",
				Sender:    "user1",
				Content:   "Hello, World!",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				RoomID:    "some-room-id",
			},
			mockErr: errors.New("database connection failed"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer mock.Close()

			expectation := mock.ExpectExec("INSERT INTO messages \\(sender, content, created_at, updated_at, room_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").
				WithArgs(tt.message.Sender, tt.message.Content, tt.message.CreatedAt, tt.message.UpdatedAt, tt.message.RoomID)

			if tt.mockErr != nil {
				expectation.WillReturnError(tt.mockErr)
			} else {
				expectation.WillReturnResult(pgxmock.NewResult("INSERT", 1))
			}

			// Create repository with mock connection
			repo := NewMessageRepository(mock)
			err = repo.CreateMessage(context.Background(), tt.message)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

// Test for CreateMessage with nil message
func TestMessageRepository_CreateMessage_NilMessage(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	repo := NewMessageRepository(mock)

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("CreateMessage with nil message should panic")
		}
	}()

	_ = repo.CreateMessage(context.Background(), nil)
}

// Test for GetMessageByID method
func TestMessageRepository_GetMessageByID(t *testing.T) {
	testTime := time.Now()

	tests := []struct {
		name      string
		messageID string
		mockRows  *pgxmock.Rows
		mockErr   error
		expected  *models.Message
		wantErr   bool
	}{
		{
			name:      "successful retrieval",
			messageID: "550e8400-e29b-41d4-a716-446655440000",
			// Fix: Match the SELECT query order: id, sender, content, created_at, updated_at, room_id
			mockRows: pgxmock.NewRows([]string{"id", "sender", "content", "created_at", "updated_at", "room_id"}).
				AddRow("550e8400-e29b-41d4-a716-446655440000", "user1", "Hello", testTime, testTime, "room-123"),
			mockErr: nil,
			expected: &models.Message{
				ID:        "550e8400-e29b-41d4-a716-446655440000",
				Sender:    "user1",
				Content:   "Hello",
				CreatedAt: testTime,
				UpdatedAt: testTime,
				RoomID:    "room-123", // Add this
			},
			wantErr: false,
		},
		{
			name:      "message not found",
			messageID: "00000000-0000-0000-0000-000000000000",
			mockRows:  nil,
			mockErr:   pgx.ErrNoRows,
			expected:  nil,
			wantErr:   false,
		},
		{
			name:      "database error",
			messageID: "550e8400-e29b-41d4-a716-446655440000",
			mockRows:  nil,
			mockErr:   errors.New("database connection failed"),
			expected:  nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer mock.Close()

			expectation := mock.ExpectQuery("SELECT \\* FROM messages where id = \\$1").
				WithArgs(tt.messageID)

			switch {
			case tt.mockErr != nil:
				expectation.WillReturnError(tt.mockErr)
			case tt.mockRows != nil:
				expectation.WillReturnRows(tt.mockRows)
			default:
				// Fix: Include all columns including room_id
				expectation.WillReturnRows(pgxmock.NewRows([]string{"id", "sender", "content", "created_at", "updated_at", "room_id"}))
			}

			repo := NewMessageRepository(mock)
			result, err := repo.GetMessageByID(context.Background(), tt.messageID)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetMessageByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.expected == nil && result != nil {
				t.Errorf("GetMessageByID() expected nil, got %v", result)
				return
			}

			if tt.expected != nil && result == nil {
				t.Errorf("GetMessageByID() expected %v, got nil", tt.expected)
				return
			}

			if tt.expected != nil && result != nil {
				if result.ID != tt.expected.ID ||
					result.Sender != tt.expected.Sender ||
					result.Content != tt.expected.Content ||
					result.RoomID != tt.expected.RoomID { // Add this check
					t.Errorf("GetMessageByID() got %v, want %v", result, tt.expected)
				}
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func BenchmarkMessageRepository_CreateMessage(b *testing.B) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		b.Fatalf("failed to create mock: %v", err)
	}
	defer mock.Close()

	message := &models.Message{
		ID:        "550e8400-e29b-41d4-a716-446655440000",
		Sender:    "user1",
		Content:   "Benchmark message",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		RoomID:    "some-room-id",
	}

	repo := NewMessageRepository(mock)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO messages \\(sender, content, created_at, updated_at, room_id\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5\\)").
			WithArgs(message.Sender, message.Content, message.CreatedAt, message.UpdatedAt, message.RoomID).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		_ = repo.CreateMessage(context.Background(), message)
	}
}
