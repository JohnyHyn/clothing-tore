package service

import (
	"clothing-store/internal/model"
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	DB *sql.DB
}

var jwtSecret = []byte("secret_key_123") // Should be in ENV

func (s *AuthService) Login(email, password string) (string, error) {
	var user model.User

	err := s.DB.QueryRow(
		"SELECT id, email, password, role FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid email or password")
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	); err != nil {
		return "", errors.New("invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func (s *AuthService) Register(user model.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if user.Role == "" {
		user.Role = "user"
	}

	_, err = s.DB.Exec(
		"INSERT INTO users (username, email, password, role) VALUES (?, ?, ?, ?)",
		user.Username, user.Email, string(hashedPassword), user.Role,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) ForgotPasswordRequest(email string) (string, error) {
	// In a real app, we would:
	// 1. Check if user exists
	// 2. Generate a random token
	// 3. Store token in DB with expiration
	// 4. Send email to user

	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&exists)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("user with this email does not exist")
	}

	// For demonstration, we'll return a dummy token "RESET123"
	// In reality, this should be a UUID or high-entropy string
	return "RESET123", nil
}

func (s *AuthService) ResetPassword(email, token, newPassword string) error {
	// In a real app, we would:
	// 1. Validate the token and email against DB
	// 2. Check if token is expired

	if token != "RESET123" {
		return errors.New("invalid or expired reset token")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.DB.Exec(
		"UPDATE users SET password = ? WHERE email = ?",
		string(hashedPassword), email,
	)
	return err
}
