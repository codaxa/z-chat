// Package postgres provides a function to establish a connection to a PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"time"
	"z-chat/internal/config"
)

// DB interface that can work with both pgx and GORM
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

// NewConnection creates a new pgx connection pool
func NewConnection() (*pgxpool.Pool, error) {
	cfg := config.New()
	dbname := cfg.DBName
	dbuser := cfg.DBUser
	dbpassword := cfg.DBPassword
	dbhost := cfg.DBHost
	dbport := cfg.DBPort

	if dbname == "" || dbuser == "" || dbpassword == "" || dbhost == "" {
		return nil, fmt.Errorf("missing required database configuration")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		url.QueryEscape(dbuser),
		url.QueryEscape(dbpassword),
		url.QueryEscape(dbhost),
		dbport,
		url.QueryEscape(dbname))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return conn, nil
}
