package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBUrl      string
	ServerPort int
}

func NewConfig() *Config {
	return &Config{
		DBUrl:      getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/conference_booking?sslmode=disable"),
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
