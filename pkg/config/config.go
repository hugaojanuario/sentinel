package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST                string
	DB_PORT                string
	DB_NAME                string
	DB_USER                string
	DB_PASSWORD            string
	DB_SSLMODE             string
	TELEGRAM_BOT_TOKEN     string
	TELEGRAM_CHAT_ID       string
	ALERT_CPU_THRESHOLD    float64
	ALERT_MEM_THRESHOLD_MB uint64
	CHECK_INTERVAL         time.Duration
}

func LoadDotEnv() *Config {
	if err := godotenv.Load("./.env"); err != nil {
		fmt.Printf("aviso: erro ao ler a .env: %v\n", err)
	}

	cpuThreshold, _ := strconv.ParseFloat(os.Getenv("ALERT_CPU_THRESHOLD"), 64)
	if cpuThreshold == 0 {
		cpuThreshold = 80
	}

	memThreshold, _ := strconv.ParseUint(os.Getenv("ALERT_MEM_THRESHOLD_MB"), 10, 64)
	if memThreshold == 0 {
		memThreshold = 500
	}

	checkInterval, _ := time.ParseDuration(os.Getenv("CHECK_INTERVAL"))
	if checkInterval == 0 {
		checkInterval = 30 * time.Second
	}

	return &Config{
		DB_HOST:                os.Getenv("DB_HOST"),
		DB_PORT:                os.Getenv("DB_PORT"),
		DB_NAME:                os.Getenv("DB_NAME"),
		DB_USER:                os.Getenv("DB_USER"),
		DB_PASSWORD:            os.Getenv("DB_PASSWORD"),
		DB_SSLMODE:             os.Getenv("DB_SSLMODE"),
		TELEGRAM_BOT_TOKEN:     os.Getenv("TELEGRAM_BOT_TOKEN"),
		TELEGRAM_CHAT_ID:       os.Getenv("TELEGRAM_CHAT_ID"),
		ALERT_CPU_THRESHOLD:    cpuThreshold,
		ALERT_MEM_THRESHOLD_MB: memThreshold,
		CHECK_INTERVAL:         checkInterval,
	}
}
