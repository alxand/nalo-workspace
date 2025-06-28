package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	if os.Getenv("DSN") == "" {
		log.Fatal("DSN environment variable is not set")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
}

func GetDSN() string {
	return os.Getenv("DSN")
}
func GetTestDSN() string {
	conn := os.Getenv("TEST_DSN")
	if conn == "" {
		conn = ":memory:"
	}
	return conn
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}
