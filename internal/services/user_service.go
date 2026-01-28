package services

import (
	"database/sql"
	"fmt"
	"jimpitan/backend/internal/database"
	"jimpitan/backend/internal/models"
	"jimpitan/backend/internal/utils"
	"time"
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

// GetAllUsers returns all active users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	rows, err := s.db.Query(
		"SELECT id, name, role, username, created_at, updated_at, last_login FROM users WHERE deleted_at IS NULL ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(&u.ID, &u.Name, &u.Role, &u.Username, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}

	return users, rows.Err()
}

// GetUserByID returns user by ID
func (s *UserService) GetUserByID(id string) (*models.User, error) {
	var u models.User
	err := s.db.QueryRow(
		"SELECT id, name, role, username, created_at, updated_at, last_login FROM users WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&u.ID, &u.Name, &u.Role, &u.Username, &u.CreatedAt, &u.UpdatedAt, &u.LastLogin)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &u, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(name, role, username, password string) (*models.User, error) {
	if name == "" || role == "" || username == "" || password == "" {
		return nil, fmt.Errorf("semua field harus diisi")
	}

	if role != "admin" && role != "petugas" {
		return nil, fmt.Errorf("role harus 'admin' atau 'petugas'")
	}

	// Check if username already exists
	var existingID string
	err := s.db.QueryRow("SELECT id FROM users WHERE username = ? AND deleted_at IS NULL", username).Scan(&existingID)
	if err != sql.ErrNoRows {
		return nil, fmt.Errorf("username sudah terdaftar")
	}

	// Get next user number
	var count int
	err = s.db.QueryRow("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get user count: %w", err)
	}

	userID := utils.GenerateUserID(count)
	passwordHash := utils.HashPassword(password)
	now := time.Now()

	_, err = s.db.Exec(
		"INSERT INTO users (id, name, role, username, password_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, name, role, username, passwordHash, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &models.User{
		ID:        userID,
		Name:      name,
		Role:      role,
		Username:  username,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateUser updates user data
func (s *UserService) UpdateUser(id, name, role, username string) error {
	if name == "" && role == "" && username == "" {
		return fmt.Errorf("setidaknya satu field harus diubah")
	}

	if role != "" && role != "admin" && role != "petugas" {
		return fmt.Errorf("role harus 'admin' atau 'petugas'")
	}

	query := "UPDATE users SET "
	args := []interface{}{}

	if name != "" {
		query += "name = ?, "
		args = append(args, name)
	}
	if role != "" {
		query += "role = ?, "
		args = append(args, role)
	}
	if username != "" {
		query += "username = ?, "
		args = append(args, username)
	}

	query += "updated_at = ? WHERE id = ? AND deleted_at IS NULL"
	args = append(args, time.Now(), id)

	_, err := s.db.Exec(query, args...)
	return err
}

// UpdatePassword updates user password
func (s *UserService) UpdatePassword(id, oldPassword, newPassword string) error {
	var passwordHash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = ? AND deleted_at IS NULL", id).Scan(&passwordHash)
	if err != nil {
		return fmt.Errorf("user tidak ditemukan")
	}

	if !utils.VerifyPassword(oldPassword, passwordHash) {
		return fmt.Errorf("password lama tidak sesuai")
	}

	newHash := utils.HashPassword(newPassword)
	_, err = s.db.Exec(
		"UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?",
		newHash, time.Now(), id,
	)
	return err
}

// DeleteUser soft deletes a user
func (s *UserService) DeleteUser(id string) error {
	_, err := s.db.Exec(
		"UPDATE users SET deleted_at = ? WHERE id = ?",
		time.Now(), id,
	)
	return err
}

// BulkDeleteUsers soft deletes multiple users
func (s *UserService) BulkDeleteUsers(ids []string) error {
	if len(ids) == 0 {
		return fmt.Errorf("tidak ada user untuk dihapus")
	}

	// Build placeholders
	query := "UPDATE users SET deleted_at = ? WHERE id IN ("
	args := []interface{}{time.Now()}
	for i := 0; i < len(ids); i++ {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args = append(args, ids[i])
	}
	query += ")"

	_, err := s.db.Exec(query, args...)
	return err
}

// GetUserActivity returns all transactions for a user
func (s *UserService) GetUserActivity(userID string) ([]models.Transaction, error) {
	rows, err := s.db.Query(
		"SELECT id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at FROM transactions WHERE user_id = ? AND deleted_at IS NULL ORDER BY timestamp DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.Timestamp, &t.CustomerID, &t.Blok, &t.Nama, &t.Nominal, &t.UserID, &t.Petugas, &t.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	return transactions, rows.Err()
}
