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
	return &Config{
		Port:       ":8080",
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     5432,
	}
}
