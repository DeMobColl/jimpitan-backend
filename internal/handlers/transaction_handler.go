package handlers

import (
	"encoding/json"
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

// DeleteTransaction soft deletes a transaction
func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := h.transactionService.DeleteTransaction(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Transaction deleted successfully", nil)
}
