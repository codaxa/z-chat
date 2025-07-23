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
		Close()
	}
}

// NewMessageRepository creates a new instance of MessageRepository with the provided database connection.
func NewMessageRepository(db interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Close()
}) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

// CreateMessage inserts a new message into the messages table.
func (m *MessageRepository) CreateMessage(ctx context.Context, msg *models.Message) error {
	query := `INSERT INTO messages (sender, receiver, content, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := m.db.Exec(ctx, query, msg.Sender, msg.Receiver, msg.Content, msg.CreatedAt, msg.UpdatedAt)
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
