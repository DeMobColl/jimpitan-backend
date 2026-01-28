package middleware

import (
	"fmt"
	"jimpitan/backend/internal/config"
	"jimpitan/backend/internal/models"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT claims
type Claims struct {
	UserID string
	Role   string
	jwt.RegisteredClaims
}

// AuthMiddleware checks JWT token validity
func AuthMiddleware(cfg *config.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := extractToken(r)
			if token == "" {
				respondError(w, http.StatusUnauthorized, "Token tidak ditemukan")
				return
			}

			claims, err := verifyToken(token, cfg.Secret)
			if err != nil {
				respondError(w, http.StatusUnauthorized, "Token tidak valid atau sudah kadaluarsa")
				return
			}

			// Store claims in context (you can use context.WithValue if needed)
			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Role", claims.Role)

			next.ServeHTTP(w, r)
		})
	}
}

// AdminOnlyMiddleware ensures user has admin role
func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		if role != "admin" {
			respondError(w, http.StatusForbidden, "Anda tidak memiliki akses ke resource ini")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func extractToken(r *http.Request) string {
	// Try Authorization header first
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
	}

	// Try query parameter as fallback
	return r.URL.Query().Get("token")
}

func verifyToken(tokenString, secret string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func respondError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := models.GenericResponse{
		Status:  "error",
		Message: message,
	}
	// Simple JSON encoding without external libraries
	fmt.Fprintf(w, `{"status":"%s","message":"%s"}`, response.Status, escapeJSON(response.Message))
}

func escapeJSON(s string) string {
	// Basic JSON escaping
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", `\n`)
	s = strings.ReplaceAll(s, "\r", `\r`)
	return s
}
