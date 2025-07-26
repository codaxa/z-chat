package postgres

import (
	"errors"
	"z-chat/internal/domain/models"

	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository is a repository that provides user storage operations using PostgreSQL database
type UserRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates and returns a new UserRepo instance with the provided database connection
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

// GetUserByUsername retrieves a user from the database by their username
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := r.db.QueryRow(ctx, query, username)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user from the database by their email
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)
	var user models.User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(ctx context.Context, u models.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, u.Username, u.Email, u.Password)
	if err != nil {
		return err
	}
	return nil
}
