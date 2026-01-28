package handlers

import (
	"encoding/json"
	"jimpitan/backend/internal/services"
	"net/http"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

// GetCustomers returns all customers
func (h *CustomerHandler) GetCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	customers, err := h.customerService.GetAllCustomers()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Customers retrieved successfully", map[string]interface{}{
		"customers": customers,
	})
}

// GetCustomerByQRHash returns customer by QR hash
func (h *CustomerHandler) GetCustomerByQRHash(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	qrHash := r.URL.Query().Get("qr_hash")
	if qrHash == "" {
		respondError(w, http.StatusBadRequest, "qr_hash parameter is required")
		return
	}

	customer, err := h.customerService.GetCustomerByQRHash(qrHash)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Customer found", customer)
}

// CreateCustomer creates a new customer
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Blok string `json:"blok"`
		Nama string `json:"nama"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	customer, err := h.customerService.CreateCustomer(req.Blok, req.Nama)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondSuccess(w, http.StatusCreated, "Customer created successfully", customer)
}

// UpdateCustomer updates customer data
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	var req struct {
		Blok string `json:"blok"`
		Nama string `json:"nama"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.customerService.UpdateCustomer(id, req.Blok, req.Nama); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Customer updated successfully", nil)
}

// DeleteCustomer soft deletes a customer
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		respondError(w, http.StatusBadRequest, "id parameter is required")
		return
	}

	if err := h.customerService.DeleteCustomer(id); err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Customer deleted successfully", nil)
}

// GetCustomerHistory returns all transactions for a customer
func (h *CustomerHandler) GetCustomerHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	customerID := r.URL.Query().Get("customer_id")
	if customerID == "" {
		respondError(w, http.StatusBadRequest, "customer_id parameter is required")
		return
	}

	transactions, err := h.customerService.GetCustomerHistory(customerID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Customer history retrieved", transactions)
}

// BulkDeleteCustomers soft deletes multiple customers
func (h *CustomerHandler) BulkDeleteCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
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
		respondError(w, http.StatusBadRequest, "ids array cannot be empty")
		return
	}

	for _, id := range req.IDs {
		if err := h.customerService.DeleteCustomer(id); err != nil {
			respondError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	respondSuccess(w, http.StatusOK, "Customers deleted successfully", map[string]int{"deleted_count": len(req.IDs)})
}
