// Package config provides the configuration for the application.
package config

import "os"

// Config holds the configuration settings for the chat server.
type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBHost     string
	DBName     string
	DBPort     int
}

// New creates a new Config instance with default values.
func New() *Config {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbName == "" {
		panic("Missing required database environment variables")
	}

	return &Config{
		Port:       ":8080",
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBHost:     dbHost,
		DBName:     dbName,
		DBPort:     5432,
	}
}
