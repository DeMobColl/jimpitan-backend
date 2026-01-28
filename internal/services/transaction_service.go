package services

import (
	"database/sql"
	"fmt"
	"jimpitan/backend/internal/database"
	"jimpitan/backend/internal/models"
	"jimpitan/backend/internal/utils"
	"time"
)

type TransactionService struct {
	db *database.DB
}

func NewTransactionService(db *database.DB) *TransactionService {
	return &TransactionService{db: db}
}

// GetAllTransactions returns all active transactions
func (s *TransactionService) GetAllTransactions() ([]models.Transaction, error) {
	rows, err := s.db.Query(
		"SELECT id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at FROM transactions WHERE deleted_at IS NULL ORDER BY timestamp DESC",
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

// GetUserTransactions returns all active transactions for a specific user
func (s *TransactionService) GetUserTransactions(userID string) ([]models.Transaction, error) {
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

// SubmitTransaction creates a new transaction
func (s *TransactionService) SubmitTransaction(customerID, blok, nama, userID, petugas string, nominal float64) (*models.Transaction, error) {
	if customerID == "" || userID == "" || nominal <= 0 {
		return nil, fmt.Errorf("data tidak lengkap atau tidak valid")
	}

	// Get next transaction number
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM transactions").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction count: %w", err)
	}

	txID := utils.GenerateTXID(count)
	now := time.Now()

	_, err = s.db.Exec(
		"INSERT INTO transactions (id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		txID, now, customerID, blok, nama, nominal, userID, petugas, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	// Update customer stats
	customerService := NewCustomerService(s.db)
	if err := customerService.UpdateCustomerStats(customerID, nominal); err != nil {
		// Log error but don't fail the transaction creation
		fmt.Printf("Warning: failed to update customer stats: %v\n", err)
	}

	return &models.Transaction{
		ID:         txID,
		Timestamp:  now,
		CustomerID: customerID,
		Blok:       blok,
		Nama:       nama,
		Nominal:    nominal,
		UserID:     userID,
		Petugas:    petugas,
		CreatedAt:  now,
	}, nil
}

// DeleteTransaction soft deletes a transaction
func (s *TransactionService) DeleteTransaction(id string) error {
	// Get transaction details for customer stats rollback
	var t models.Transaction
	err := s.db.QueryRow(
		"SELECT id, customer_id, nominal FROM transactions WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&t.ID, &t.CustomerID, &t.Nominal)
	if err != nil {
		return fmt.Errorf("transaction tidak ditemukan")
	}

	// Soft delete transaction
	_, err = s.db.Exec(
		"UPDATE transactions SET deleted_at = ? WHERE id = ?",
		time.Now(), id,
	)
	if err != nil {
		return err
	}

	// Rollback customer stats
	_, err = s.db.Exec(
		"UPDATE customers SET total_setoran = total_setoran - ?, updated_at = ? WHERE id = ?",
		t.Nominal, time.Now(), t.CustomerID,
	)
	return err
}

// DeleteTransactionWithValidation soft deletes a transaction after validating user permission
func (s *TransactionService) DeleteTransactionWithValidation(id, userID, userRole string) error {
	// Get transaction details
	var t models.Transaction
	err := s.db.QueryRow(
		"SELECT id, customer_id, nominal, user_id FROM transactions WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&t.ID, &t.CustomerID, &t.Nominal, &t.UserID)
	if err != nil {
		return fmt.Errorf("transaksi tidak ditemukan")
	}

	// Validation: petugas hanya bisa hapus transaksi miliknya sendiri
	if userRole != "admin" && t.UserID != userID {
		return fmt.Errorf("anda hanya dapat menghapus transaksi milik anda sendiri")
	}

	// Soft delete transaction
	_, err = s.db.Exec(
		"UPDATE transactions SET deleted_at = ? WHERE id = ?",
		time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("gagal menghapus transaksi: %w", err)
	}

	// Rollback customer stats
	_, err = s.db.Exec(
		"UPDATE customers SET total_setoran = total_setoran - ?, updated_at = ? WHERE id = ?",
		t.Nominal, time.Now(), t.CustomerID,
	)
	if err != nil {
		return fmt.Errorf("gagal memperbarui statistik customer: %w", err)
	}

	return nil
}

// BulkDeleteTransactions soft deletes multiple transactions (admin only)
// Returns count of deleted transactions and slice of errors
func (s *TransactionService) BulkDeleteTransactions(ids []string) (int, []map[string]string) {
	var deleted int
	var errors []map[string]string

	for _, id := range ids {
		// Get transaction details
		var t models.Transaction
		err := s.db.QueryRow(
			"SELECT id, customer_id, nominal FROM transactions WHERE id = ? AND deleted_at IS NULL",
			id,
		).Scan(&t.ID, &t.CustomerID, &t.Nominal)
		if err != nil {
			errors = append(errors, map[string]string{
				"id":    id,
				"error": "transaksi tidak ditemukan",
			})
			continue
		}

		// Soft delete transaction
		_, err = s.db.Exec(
			"UPDATE transactions SET deleted_at = ? WHERE id = ?",
			time.Now(), id,
		)
		if err != nil {
			errors = append(errors, map[string]string{
				"id":    id,
				"error": "gagal menghapus transaksi",
			})
			continue
		}

		// Rollback customer stats
		_, err = s.db.Exec(
			"UPDATE customers SET total_setoran = total_setoran - ?, updated_at = ? WHERE id = ?",
			t.Nominal, time.Now(), t.CustomerID,
		)
		if err != nil {
			errors = append(errors, map[string]string{
				"id":    id,
				"error": "gagal memperbarui statistik customer",
			})
			continue
		}

		deleted++
	}

	return deleted, errors
}

// GetTransactionByID returns transaction by ID
func (s *TransactionService) GetTransactionByID(id string) (*models.Transaction, error) {
	var t models.Transaction
	err := s.db.QueryRow(
		"SELECT id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at FROM transactions WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&t.ID, &t.Timestamp, &t.CustomerID, &t.Blok, &t.Nama, &t.Nominal, &t.UserID, &t.Petugas, &t.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("transaction tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &t, nil
}
