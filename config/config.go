package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT   int
	DB     string
	SECRET string
}

var CONFIG Config

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	CONFIG = Config{
		PORT:   getIntEnv("PORT", 8000),
		DB:     buildDataSourceName(),
		SECRET: os.Getenv("SECRET"),
	}
}

func buildDataSourceName() string {
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	return dsn
}

func getIntEnv(key string, defaultValue int) int {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}

	return defaultValue
}
