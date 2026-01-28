package services

import (
	"database/sql"
	"fmt"
	"jimpitan/backend/internal/database"
	"jimpitan/backend/internal/models"
	"jimpitan/backend/internal/utils"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	db        *database.DB
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(db *database.DB, jwtSecret string, expiryHours int) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
		jwtExpiry: expiryHours,
	}
}

// Login authenticates user and returns token
func (s *AuthService) Login(username, password string) (*models.LoginResponse, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username dan password harus diisi")
	}

	var user models.User
	var dbPasswordHash string

	err := s.db.QueryRow(
		"SELECT id, name, role, username, password_hash FROM users WHERE username = ? AND deleted_at IS NULL",
		username,
	).Scan(&user.ID, &user.Name, &user.Role, &user.Username, &dbPasswordHash)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("username atau password salah")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if !utils.VerifyPassword(password, dbPasswordHash) {
		return nil, fmt.Errorf("username atau password salah")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.jwtExpiry)).Unix(),
		"iat":     time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update user with new token in database
	now := time.Now()
	tokenExpiry := now.Add(time.Hour * time.Duration(s.jwtExpiry))

	_, err = s.db.Exec(
		"UPDATE users SET token = ?, token_expiry = ?, last_login = ?, updated_at = ? WHERE id = ?",
		tokenString, tokenExpiry, now, now, user.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user token: %w", err)
	}

	return &models.LoginResponse{
		ID:          user.ID,
		Name:        user.Name,
		Role:        user.Role,
		Username:    user.Username,
		Token:       tokenString,
		TokenExpiry: tokenExpiry,
		LastLogin:   now,
	}, nil
}

// VerifyToken verifies if token is still valid
func (s *AuthService) VerifyToken(token string) (*models.User, error) {
	if token == "" {
		return nil, fmt.Errorf("token tidak ditemukan")
	}

	var user models.User
	var tokenExpiry time.Time

	err := s.db.QueryRow(
		"SELECT id, name, role, username, token, token_expiry FROM users WHERE token = ? AND deleted_at IS NULL",
		token,
	).Scan(&user.ID, &user.Name, &user.Role, &user.Username, &user.Token, &tokenExpiry)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("token tidak valid")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Check token expiry
	if time.Now().After(tokenExpiry) {
		return nil, fmt.Errorf("token sudah kadaluarsa")
	}

	user.TokenExpiry = tokenExpiry
	return &user, nil
}

// Logout invalidates user's token
func (s *AuthService) Logout(userID string) error {
	_, err := s.db.Exec(
		"UPDATE users SET token = NULL, token_expiry = NULL, updated_at = ? WHERE id = ?",
		time.Now(), userID,
	)
	return err
}
