package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	ServerPort        string
	ServerMode        string
	DBHost            string
	DBUser            string
	DBPassword        string
	DBName            string
	EncryptionKey     string
	EncryptionSignKey string
}

var config Configuration

func init() {
	godotenv.Load()

	config = Configuration{
		ServerPort:        getEnv("SERVER_PORT", ""),
		ServerMode:        getEnv("GIN_MODE", ""),
		DBHost:            getEnv("DB_HOST", ""),
		DBUser:            getEnv("DB_USER", ""),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", ""),
		EncryptionKey:     getEnv("ENCRYPTION_KEY", ""),
		EncryptionSignKey: getEnv("ENCRYPTION_SIGN_KEY", ""),
	}
}

func GetConfig() Configuration {
	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
