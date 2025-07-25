// Package config provides the configuration for the application.
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds the configuration settings for the chat server.
type Config struct {
	Port          string
	DBUser        string
	DBPassword    string
	DBHost        string
	DBName        string
	DBPort        int
	JWTSecret     string
	TokenDuration time.Duration
	DBUrl         string
}

// New creates a new Config instance with default values.
func New() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Use default values for development if environment variables are not set
	if dbUser == "" {
		dbUser = "postgres"
		log.Println("Warning: Using default database user. Set DB_USER environment variable in production.")
	}
	if dbPassword == "" {
		dbPassword = "postgres"
		log.Println("Warning: Using default database password. Set DB_PASSWORD environment variable in production.")
	}
	if dbHost == "" {
		dbHost = "localhost"
		log.Println("Warning: Using default database host. Set DB_HOST environment variable in production.")
	}
	if dbName == "" {
		dbName = "z-chat"
		log.Println("Warning: Using default database name. Set DB_NAME environment variable in production.")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-secret-for-development-only"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}

	dbURL := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, 5432, dbUser, dbPassword, dbName)

	return &Config{
		Port:          ":8080",
		DBUser:        dbUser,
		DBPassword:    dbPassword,
		DBHost:        dbHost,
		DBName:        dbName,
		DBPort:        5432,
		JWTSecret:     jwtSecret,
		TokenDuration: 24 * time.Hour,
		DBUrl:         dbURL,
	}
}
