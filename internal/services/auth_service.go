// Package services provides business logic implementations for the application.
package services

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
	"z-chat/internal/domain/models"
	"z-chat/internal/domain/repository"
)

// AuthService handles authentication operations.
type AuthService struct {
	userRepository repository.UserRepository
	jwtSecret      []byte
	tokenDuration  time.Duration
}

// NewAuthService creates and returns a new instance of AuthService.
func NewAuthService(userRepository repository.UserRepository, jwtSecret string, tokenDuration time.Duration) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		jwtSecret:      []byte(jwtSecret),
		tokenDuration:  tokenDuration,
	}
}

// Register creates a new user account with the provided username and password.
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	return s.userRepository.CreateUser(ctx, user)
}

// Login authenticates a user with the provided username and password.
func (s *AuthService) Login(ctx context.Context, identifier, password string) (string, error) {
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

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

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

// ValidateToken checks the validity of a JWT token and returns the claims if valid.
func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if token == nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("failed to parse token claims")
	}

	err = claims.Valid()
	if err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %w", err)
	}

	return &claims, nil
}
