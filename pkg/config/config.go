package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
	DB_SSLMODE  string
}

func LoadDotEnv() *Config {
	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Errorf("error : erro ao ler a .env: %w", err)
	}

	return &Config{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_SSLMODE:  os.Getenv("DB_SSLMODE"),
	}
}
