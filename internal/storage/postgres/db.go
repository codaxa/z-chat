// Package postgres provides a function to establish a connection to a PostgreSQL database.
package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
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

	var connectionString = "postgres://" + dbuser + ":" + dbpassword + "@" + dbhost + ":" + strconv.Itoa(dbport) + "/" + dbname

	conn, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
