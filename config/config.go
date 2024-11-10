package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RetryCount int
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return Config{
		RetryCount: getInt("RETRY_COUNT"),
	}
}

func getInt(value string) int {
	intValue, err := strconv.Atoi(os.Getenv(value))
	if err != nil {
		log.Printf("Error converting %s to int: %v.", value, err)
	}
	return intValue
}
