// Package postgres provides a function to establish a connection to a PostgreSQL database.
package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
	"time"
	"z-chat/internal/config"
	"z-chat/internal/domain/models"
)

// DB interface that can work with both pgx and GORM
type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Close()
}

// GormDB interface for GORM operations
type GormDB interface {
	Create(value interface{}) *gorm.DB
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
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

// NewGormConnection creates a new GORM connection for migrations
func NewGormConnection() (*gorm.DB, error) {
	cfg := config.New()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// AutoMigrate runs GORM auto-migration for development
func AutoMigrate() error {
	db, err := NewGormConnection()
	if err != nil {
		return err
	}

	return db.AutoMigrate(
		&models.User{},
		&models.Message{},
	)
}
