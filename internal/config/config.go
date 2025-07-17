// Package config provides the configuration for the application.
package config

// Config holds the configuration settings for the chat server.
type Config struct {
	Port string
}

// New returns a pointer to a Config initialized with the default server port ":8080".
func New() *Config {
	return &Config{
		Port: ":8080",
	}
}
