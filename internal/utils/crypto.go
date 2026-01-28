package utils

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// VerifyPassword verifies password against hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// GenerateToken generates a 32-character random token
func GenerateToken() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	token := make([]byte, 32)
	for i := range token {
		token[i] = charset[rand.Intn(len(charset))]
	}
	return string(token)
}

// GenerateQRHash generates a 10-character QR hash
func GenerateQRHash(customerID string) string {
	saltedString := "Jimpitan" + customerID
	hash := sha256.Sum256([]byte(saltedString))
	// Convert to hex and take first 10 characters
	hashHex := fmt.Sprintf("%x", hash)
	if len(hashHex) > 10 {
		hashHex = hashHex[:10]
	}
	return hashHex
}

// GenerateUserID generates a user ID in format USR-XXX
func GenerateUserID(count int) string {
	return fmt.Sprintf("USR-%03d", count+1)
}

// GenerateCustomerID generates a customer ID in format CUST-XXX
func GenerateCustomerID(count int) string {
	return fmt.Sprintf("CUST-%03d", count+1)
}

// GenerateTXID generates a transaction ID
func GenerateTXID(count int) string {
	return fmt.Sprintf("%04d", count+1)
}

// GetCurrentTimestamp returns current timestamp in ISO format
func GetCurrentTimestamp() time.Time {
	return time.Now().UTC()
}
