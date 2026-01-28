package handlers

import (
	"encoding/json"
	"fmt"
	"jimpitan/backend/internal/services"
	"net/http"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
}

func NewTransactionHandler(transactionService *services.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// GetHistory returns all transactions
func (h *TransactionHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	transactions, err := h.transactionService.GetAllTransactions()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Transactions retrieved successfully", transactions)
}

// GetMyHistory returns transactions for the current user
func (h *TransactionHandler) GetMyHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "User information not found")
		return
	}

	transactions, err := h.transactionService.GetUserTransactions(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "User transactions retrieved successfully", transactions)
}

// SubmitTransaction creates a new transaction
func (h *TransactionHandler) SubmitTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		CustomerID string  `json:"customer_id"`
		Blok       string  `json:"blok"`
		Nama       string  `json:"nama"`
		Nominal    float64 `json:"nominal"`
		UserID     string  `json:"user_id"`
		Petugas    string  `json:"petugas"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction, err := h.transactionService.SubmitTransaction(
		req.CustomerID,
		req.Blok,
		req.Nama,
		req.UserID,
		req.Petugas,
		req.Nominal,
	)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, "Transaksi berhasil ditambahkan", transaction)
}

// DeleteTransaction soft deletes a single transaction with validation
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	userRole := r.Header.Get("X-User-Role")

	if userID == "" || userRole == "" {
		respondError(w, http.StatusUnauthorized, "User information not found")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := h.transactionService.DeleteTransactionWithValidation(id, userID, userRole); err != nil {
		respondError(w, http.StatusForbidden, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Transaksi berhasil dihapus", nil)
}

// BulkDeleteTransactions soft deletes multiple transactions (admin only)
func (h *TransactionHandler) BulkDeleteTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userRole := r.Header.Get("X-User-Role")
	if userRole != "admin" {
		respondError(w, http.StatusForbidden, "Hanya admin yang dapat menghapus transaksi massal")
		return
	}

	var req struct {
		IDs []string `json:"ids"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if len(req.IDs) == 0 {
		respondError(w, http.StatusBadRequest, "ids array is required and cannot be empty")
		return
	}

	deleted, errors := h.transactionService.BulkDeleteTransactions(req.IDs)

	if len(errors) > 0 && deleted == 0 {
		respondError(w, http.StatusInternalServerError, "Semua transaksi gagal dihapus")
		return
	}

	response := map[string]interface{}{
		"deleted": deleted,
	}
	if len(errors) > 0 {
		response["errors"] = errors
	}

	respondSuccess(w, http.StatusOK, fmt.Sprintf("%d transaksi berhasil dihapus", deleted), response)
}
