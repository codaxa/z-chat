// Package config provides the configuration for the application.
package config

// Config holds the configuration settings for the chat server.
type Config struct {
	Port string
}

// New creates a new Config instance with default values.
func New() *Config {
	return &Config{
		Port: ":8080",
	}
}
