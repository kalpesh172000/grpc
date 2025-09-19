package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
	LogFile     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://username:password@localhost/dbname?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		Port:        getEnv("PORT", "8080"),
		LogFile:     getEnv("LOG_FILE", "server.log"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
