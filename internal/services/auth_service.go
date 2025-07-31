// Package services provides business logic implementations for the application.
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
	userRepository repository.UserRepository
	jwtSecret      []byte
	tokenDuration  time.Duration
}

// NewAuthService creates and returns a new instance of AuthService with the provided
// user repository, JWT secret, and token duration.
func NewAuthService(userRepository repository.UserRepository, jwtSecret string, tokenDuration time.Duration) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      []byte(jwtSecret),
		tokenDuration:  tokenDuration,
	}
}

// Register creates a new user account with the provided username and password.
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	// 1. Check if username already exists
	existingUser, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return fmt.Errorf("failed to check username availability: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("username already exists")
	}

	existingUser, err = s.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to check email availability: %w", err)
	}
	if existingUser != nil {
		return fmt.Errorf("email already exists")
	}

	// 2. Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 3. Save the new user to the database
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepository.CreateUser(ctx, user)
}

// Login authenticates a user with the provided username and password.
func (s *AuthService) Login(ctx context.Context, identifier, password string) (string, error) {
	// 1. Get user from repository
	userByUsername, err := s.userRepository.GetUserByUsername(ctx, identifier)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user: %w", err)
	}

	userByEmail, err := s.userRepository.GetUserByEmail(ctx, identifier)
	if err != nil {
		return "", fmt.Errorf("failed to fetch user by email: %w", err)
	}
	if userByUsername == nil && userByEmail == nil {
		return "", fmt.Errorf("invalid username or email")
	}
	var user *models.User
	if userByUsername != nil {
		user = userByUsername
	} else {
		user = userByEmail
	}

	// 2. Compare password with stored hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
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
