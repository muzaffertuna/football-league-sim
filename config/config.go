package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConnectionString string
	ServerAddress      string
}

func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using default environment variables: %v", err)
	}

	cfg := Config{
		DBConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		ServerAddress:      os.Getenv("SERVER_ADDRESS"),
	}

	if cfg.DBConnectionString == "" {
		log.Fatal("DB_CONNECTION_STRING is required")
	}
	if cfg.ServerAddress == "" {
		log.Fatal("SERVER_ADDRESS is required")
	}

	return cfg
}
