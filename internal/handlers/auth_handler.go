package handlers

import (
	"encoding/json"
	"jimpitan/backend/internal/models"
	"jimpitan/backend/internal/services"
	"net/http"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login handles user login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	loginResp, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Login berhasil", loginResp)
}

// VerifyToken verifies if token is valid
func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		respondError(w, http.StatusBadRequest, "Token parameter is required")
		return
	}

	user, err := h.authService.VerifyToken(token)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Token valid", map[string]interface{}{
		"id":            user.ID,
		"name":          user.Name,
		"role":          user.Role,
		"username":      user.Username,
		"token":         user.Token,
		"token_expiry":  user.TokenExpiry,
	})
}

// Logout handles user logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "User ID not found in token")
		return
	}

	if err := h.authService.Logout(userID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to logout: "+err.Error())
		return
	}

	respondSuccess(w, http.StatusOK, "Logout berhasil", nil)
}

// Health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, http.StatusOK, "Jimpitan App API Active", map[string]interface{}{
		"status":   "success",
		"version":  "2.0",
		"features": []string{"Login", "CRUD User", "History", "Password Hashing"},
	})
}

// Helper functions
func respondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.GenericResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
	json.NewEncoder(w).Encode(response)
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.GenericResponse{
		Status:  "error",
		Message: message,
	}
	json.NewEncoder(w).Encode(response)
}
