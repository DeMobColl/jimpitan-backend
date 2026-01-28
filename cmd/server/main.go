package main

import (
	"fmt"
	"jimpitan/backend/internal/config"
	"jimpitan/backend/internal/database"
	"jimpitan/backend/internal/handlers"
	"jimpitan/backend/internal/middleware"
	"jimpitan/backend/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewDB(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize services
	authService := services.NewAuthService(db, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	userService := services.NewUserService(db)
	customerService := services.NewCustomerService(db)
	transactionService := services.NewTransactionService(db)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/", handlers.HealthCheck).Methods(http.MethodGet)
	router.HandleFunc("/api/health", handlers.HealthCheck).Methods(http.MethodGet)

	// Auth endpoints
	router.HandleFunc("/api/login", authHandler.Login).Methods(http.MethodPost)
	router.HandleFunc("/api/verifyToken", authHandler.VerifyToken).Methods(http.MethodGet)
	router.HandleFunc("/api/logout", authHandler.Logout).Methods(http.MethodPost)

	// User endpoints (protected)
	userRoutes := router.PathPrefix("/api/users").Subrouter()
	userRoutes.Use(middleware.AuthMiddleware(&cfg.JWT))
	userRoutes.HandleFunc("", userHandler.GetUsers).Methods(http.MethodGet)
	userRoutes.HandleFunc("", userHandler.CreateUser).Methods(http.MethodPost)
	userRoutes.HandleFunc("", userHandler.UpdateUser).Methods(http.MethodPut)
	userRoutes.HandleFunc("", userHandler.DeleteUser).Methods(http.MethodDelete)
	userRoutes.HandleFunc("/activity", userHandler.GetUserActivity).Methods(http.MethodGet)
	userRoutes.HandleFunc("/bulk-delete", userHandler.BulkDeleteUsers).Methods(http.MethodPost)
	userRoutes.HandleFunc("/password", userHandler.UpdatePassword).Methods(http.MethodPost)

	// Customer endpoints (protected)
	customerRoutes := router.PathPrefix("/api/customers").Subrouter()
	customerRoutes.Use(middleware.AuthMiddleware(&cfg.JWT))
	customerRoutes.HandleFunc("", customerHandler.GetCustomers).Methods(http.MethodGet)
	customerRoutes.HandleFunc("", customerHandler.CreateCustomer).Methods(http.MethodPost)
	customerRoutes.HandleFunc("", customerHandler.UpdateCustomer).Methods(http.MethodPut)
	customerRoutes.HandleFunc("", customerHandler.DeleteCustomer).Methods(http.MethodDelete)
	customerRoutes.HandleFunc("/qr", customerHandler.GetCustomerByQRHash).Methods(http.MethodGet)
	customerRoutes.HandleFunc("/history", customerHandler.GetCustomerHistory).Methods(http.MethodGet)
	customerRoutes.HandleFunc("/bulk-delete", customerHandler.BulkDeleteCustomers).Methods(http.MethodPost)

	// Transaction endpoints (protected)
	transactionRoutes := router.PathPrefix("/api/transactions").Subrouter()
	transactionRoutes.Use(middleware.AuthMiddleware(&cfg.JWT))
	transactionRoutes.HandleFunc("", transactionHandler.GetHistory).Methods(http.MethodGet)
	transactionRoutes.HandleFunc("", transactionHandler.SubmitTransaction).Methods(http.MethodPost)
	transactionRoutes.HandleFunc("", transactionHandler.DeleteTransaction).Methods(http.MethodDelete)

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: cfg.CORS.AllowedOrigins,
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"*"},
		MaxAge:         3600,
	})

	handler := c.Handler(router)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("🚀 Jimpitan Backend Server starting on %s (ENV: %s)", addr, cfg.Server.Env)
	log.Printf("📊 Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)

	if err := http.ListenAndServe(addr, handler); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func init() {
	// Load .env file if exists
	if _, err := os.Stat(".env"); err == nil {
		// Simple .env loader - in production use godotenv package
		loadEnvFile(".env")
	}
}

func loadEnvFile(filename string) {
	// This is a placeholder - in production, use github.com/joho/godotenv
	// or implement proper env file loading
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	// Simple parsing (not production ready)
	lines := string(data)
	// Parse env variables from file
	_ = lines
}
