package postgres

import (
	"database/sql"
	"errors"
	"z-chat/internal/domain/models"

	"golang.org/x/net/context"
)

// UserRepo is a repository that provides user storage operations using PostgreSQL database
type UserRepo struct {
	db *sql.DB
}

// NewUserRepo creates and returns a new UserRepo instance with the provided database connection
func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

// GetByUsername retrieves a user from the database by their username
// Returns nil, nil if the user is not found
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, Hashed_password FROM users WHERE username = $1`
	row := r.db.QueryRowContext(ctx, query, username)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.HashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err

	}
	return &user, nil
}

// Add inserts a new user into the database
// Returns an error if the operation fails
func (r *UserRepo) Add(ctx context.Context, u models.User) error {
	query := `INSERT INTO users (username, Hashed_password) VALUES ($1, $2)`
	_, err := r.db.ExecContext(ctx, query, u.Username, u.HashedPassword)
	if err != nil {
		return err
	}
	return nil
}
