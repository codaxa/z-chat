package postgres

import (
	"context"
	"z-chat/internal/domain/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// MessageRepository provides methods to interact with the messages table in the database.
type MessageRepository struct {
	db interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		Close()
	}
}

// NewMessageRepository creates a new instance of MessageRepository with the provided database connection.
func NewMessageRepository(db interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

// CreateMessage inserts a new message into the messages table.
func (m *MessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	query := `INSERT INTO messages (sender, receiver, content, created_at, updated_at, room_id) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := m.db.Exec(ctx, query, msg.Sender, msg.Receiver, msg.Content, msg.CreatedAt, msg.UpdatedAt, msg.RoomID)
	return err
}

// GetMessageByID retrieves a message by its ID from the messages table.
func (m *MessageRepository) GetMessageByID(ctx context.Context, id string) (*models.Message, error) {
	query := `SELECT * FROM messages where id = $1`
	message := m.db.QueryRow(ctx, query, id)

	var msg models.Message

	err := message.Scan(&msg.ID, &msg.Sender, &msg.Receiver, &msg.Content, &msg.CreatedAt, &msg.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &msg, nil
}

// GetMessagesByRoom retrieves all messages for a specific room from the messages table.
// It returns the messages sorted by creation time in ascending order.
//
// Parameters:
//   - ctx: The context for the database operation
//   - roomID: The unique identifier of the room
//
// Returns:
//   - []*models.Message: A slice of message pointers for the specified room
//   - error: Any error encountered during the query execution
func (m *MessageRepository) GetMessagesByRoom(ctx context.Context, roomID string, limit, offset int) ([]*models.Message, error) {
	query := `
		SELECT id, sender, receiver, content, created_at, updated_at, room_id
		FROM messages
		WHERE room_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := m.db.Query(ctx, query, roomID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := new(models.Message)
		if err := rows.Scan(
			&msg.ID, &msg.Sender, &msg.Receiver, &msg.Content,
			&msg.CreatedAt, &msg.UpdatedAt, &msg.RoomID,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
