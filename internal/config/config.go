package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv        string
	AppPort       string
	DatabaseURL   string
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration
	CORSOrigin    string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	accessTTL, err := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "24h"))
	if err != nil || accessTTL <= 0 {
		accessTTL = 24 * time.Hour
	}
	refreshTTL, err := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))
	if err != nil || refreshTTL <= 0 {
		refreshTTL = 7 * 24 * time.Hour
	}

	return &Config{
		AppEnv:        getEnv("APP_ENV", "development"),
		AppPort:       getEnv("APP_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		JWTSecret:     getEnv("JWT_SECRET", "default-secret"),
		JWTAccessTTL:  accessTTL,
		JWTRefreshTTL: refreshTTL,
		CORSOrigin:    getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
