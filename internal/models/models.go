package models

import "time"

// User represents a system user (admin or petugas)
type User struct {
	ID           string     `json:"id"`           // USR-001
	Name         string     `json:"name"`
	Role         string     `json:"role"`         // admin or petugas
	Username     string     `json:"username"`
	PasswordHash string     `json:"-"`            // Never expose password hash
	Token        string     `json:"token,omitempty"`
	TokenExpiry  *time.Time `json:"token_expiry,omitempty"`
	LastLogin    *time.Time `json:"last_login,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// Customer represents a jimpitan member
type Customer struct {
	ID               string    `json:"id"`           // CUST-001
	Blok             string    `json:"blok"`         // Block/ID number
	Nama             string    `json:"nama"`         // Full name
	QRHash           string    `json:"qr_hash"`      // 10-char QR identifier
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	TotalSetoran     float64   `json:"total_setoran"`
	LastTransaction  *time.Time `json:"last_transaction,omitempty"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty"`
}

// Transaction represents a jimpitan deposit
type Transaction struct {
	ID           string    `json:"id"`           // TXID
	Timestamp    time.Time `json:"timestamp"`
	CustomerID   string    `json:"customer_id"`  // Reference to Customer
	Blok         string    `json:"blok"`
	Nama         string    `json:"nama"`
	Nominal      float64   `json:"nominal"`
	UserID       string    `json:"user_id"`      // Reference to User
	Petugas      string    `json:"petugas"`      // Staff name
	CreatedAt    time.Time `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// Config represents system configuration
type Config struct {
	ID                  string    `json:"id"`
	PetugasWebLoginEnabled bool   `json:"petugas_web_login_enabled"`
	MobileAppVersion    string    `json:"mobile_app_version"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// Session represents an active user session
type Session struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents successful login response
type LoginResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Role        string     `json:"role"`
	Username    string     `json:"username"`
	Token       string     `json:"token"`
	TokenExpiry *time.Time `json:"token_expiry"`
	LastLogin   *time.Time `json:"last_login"`
}

// GenericResponse represents standard API response
type GenericResponse struct {
	Status  string      `json:"status"`  // success or error
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
