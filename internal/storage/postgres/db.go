// Package postgres provides a function to establish a connection to a PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"time"
	"z-chat/internal/config"
)

// NewConnection creates a new connection pool to the PostgreSQL database using configuration settings.
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
		dbhost,
		dbport,
		dbname)

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
