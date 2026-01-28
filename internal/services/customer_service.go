package services

import (
	"database/sql"
	"fmt"
	"jimpitan/backend/internal/database"
	"jimpitan/backend/internal/models"
	"jimpitan/backend/internal/utils"
	"time"
)

type CustomerService struct {
	db *database.DB
}

func NewCustomerService(db *database.DB) *CustomerService {
	return &CustomerService{db: db}
}

// GetAllCustomers returns all active customers
func (s *CustomerService) GetAllCustomers() ([]models.Customer, error) {
	rows, err := s.db.Query(
		"SELECT id, blok, nama, qr_hash, created_at, updated_at, total_setoran, last_transaction FROM customers WHERE deleted_at IS NULL ORDER BY id",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query customers: %w", err)
	}
	defer rows.Close()

	var customers []models.Customer
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(&c.ID, &c.Blok, &c.Nama, &c.QRHash, &c.CreatedAt, &c.UpdatedAt, &c.TotalSetoran, &c.LastTransaction)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, c)
	}

	return customers, rows.Err()
}

// GetCustomerByID returns customer by ID
func (s *CustomerService) GetCustomerByID(id string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow(
		"SELECT id, blok, nama, qr_hash, created_at, updated_at, total_setoran, last_transaction FROM customers WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&c.ID, &c.Blok, &c.Nama, &c.QRHash, &c.CreatedAt, &c.UpdatedAt, &c.TotalSetoran, &c.LastTransaction)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &c, nil
}

// GetCustomerByQRHash returns customer by QR hash
func (s *CustomerService) GetCustomerByQRHash(qrHash string) (*models.Customer, error) {
	var c models.Customer
	err := s.db.QueryRow(
		"SELECT id, blok, nama, qr_hash, created_at, updated_at, total_setoran, last_transaction FROM customers WHERE qr_hash = ? AND deleted_at IS NULL",
		qrHash,
	).Scan(&c.ID, &c.Blok, &c.Nama, &c.QRHash, &c.CreatedAt, &c.UpdatedAt, &c.TotalSetoran, &c.LastTransaction)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer tidak ditemukan")
	}
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &c, nil
}

// CreateCustomer creates a new customer
func (s *CustomerService) CreateCustomer(blok, nama string) (*models.Customer, error) {
	if blok == "" || nama == "" {
		return nil, fmt.Errorf("blok dan nama harus diisi")
	}

	// Get next customer number
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&count)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer count: %w", err)
	}

	customerID := utils.GenerateCustomerID(count)
	qrHash := utils.GenerateQRHash(customerID)
	now := time.Now()

	_, err = s.db.Exec(
		"INSERT INTO customers (id, blok, nama, qr_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
		customerID, blok, nama, qrHash, now, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	return &models.Customer{
		ID:        customerID,
		Blok:      blok,
		Nama:      nama,
		QRHash:    qrHash,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// UpdateCustomer updates customer data
func (s *CustomerService) UpdateCustomer(id, blok, nama string) error {
	_, err := s.db.Exec(
		"UPDATE customers SET blok = ?, nama = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		blok, nama, time.Now(), id,
	)
	return err
}

// DeleteCustomer soft deletes a customer
func (s *CustomerService) DeleteCustomer(id string) error {
	_, err := s.db.Exec(
		"UPDATE customers SET deleted_at = ? WHERE id = ?",
		time.Now(), id,
	)
	return err
}

// GetCustomerHistory returns all transactions for a customer
func (s *CustomerService) GetCustomerHistory(customerID string) ([]models.Transaction, error) {
	rows, err := s.db.Query(
		"SELECT id, timestamp, customer_id, blok, nama, nominal, user_id, petugas, created_at FROM transactions WHERE customer_id = ? AND deleted_at IS NULL ORDER BY timestamp DESC",
		customerID,
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

// UpdateCustomerStats updates customer's total setoran and last transaction
func (s *CustomerService) UpdateCustomerStats(customerID string, amount float64) error {
	_, err := s.db.Exec(
		"UPDATE customers SET total_setoran = total_setoran + ?, last_transaction = ?, updated_at = ? WHERE id = ?",
		amount, time.Now(), time.Now(), customerID,
	)
	return err
}
