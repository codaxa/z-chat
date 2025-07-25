// Package services provides business logic implementations for the application.
// It contains service layers that coordinate operations between repositories
// and domain models while implementing various application features.
package services

import (
	"context"
	"fmt"
	"time"
	"z-chat/internal/domain/models"
	"z-chat/internal/domain/repository"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication operations such as user registration,
// login, and JWT token generation/validation.
type AuthService struct {
	userRepo      repository.UserRepo
	jwtSecret     []byte
	tokenDuration time.Duration
}

// NewAuthService creates and returns a new instance of AuthService with the provided
// user repository, JWT secret, and token duration.
func NewAuthService(userRepo repository.UserRepo, jwtSercet string, tokenDuration time.Duration) *AuthService {
	return &AuthService{
		userRepo:      userRepo,
		jwtSecret:     []byte(jwtSercet),
		tokenDuration: tokenDuration,
	}
}

// Register creates a new user account with the provided username and password.
// It checks if the username is available, hashes the password, and stores the user in the repository.
// Returns an error if the username already exists or if there are any issues with the registration process.
func (s *AuthService) Register(ctx context.Context, username, password string) error {
	// 1. Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to check username availability: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("username already exists")
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Save the new user to the database
	user := models.User{
		Username:       username,
		HashedPassword: string(hashedPassword),
	}

	return s.userRepo.Add(ctx, user)
}

// Login authenticates a user with the provided username and password.
// It retrieves the user from the repository, verifies the password hash,
// and generates a JWT token upon successful authentication.
// Returns the JWT token as a string and nil error on success, or an empty string and error on failure.
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// 1. Get user from repository
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// 2. Compare password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	// 3. Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      now.Add(s.tokenDuration).Unix(),
		"iat":      now.Unix(),
		"iss":      "z-chat",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
