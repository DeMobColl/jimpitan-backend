package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port int
	Env  string
}

type JWTConfig struct {
	Secret       string
	ExpiryHours  int
}

type CORSConfig struct {
	AllowedOrigins []string
}

func Load() *Config {
	// Load .env file (ignore error if not found)
	_ = godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	serverPort, _ := strconv.Atoi(getEnv("PORT", "8080"))
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "168"))

	corsOrigins := []string{
		"http://localhost:3000",
		"http://localhost:5173",
	}
	if origins := getEnv("CORS_ALLOWED_ORIGINS", ""); origins != "" {
		// Parse from comma-separated string if needed
		corsOrigins = append(corsOrigins, origins)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "jimpitan"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "jimpitan"),
		},
		Server: ServerConfig{
			Port: serverPort,
			Env:  getEnv("ENV", "development"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpiryHours: jwtExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: corsOrigins,
		},
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
	)
}
