package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBUser     string
	DBPassword string
	DBAddress  string
	DBName     string
	JWTSecret  string
}

var Envs = InitializeConfig()

func InitializeConfig() Config {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		fmt.Println("Error:", err)
	}

	return Config{
		Port:       getEnv("PORT", "5432"),
		DBUser:     getEnv("DB_USER", "user_2"),
		DBPassword: getEnv("DB_PASSWORD", "test"),
		DBName:     getEnv("DB_NAME", "norvista"),
		JWTSecret:  getEnv("JWT_SECRET", "c757b8c7cacc1d63b3d37a5688eaef1809687c1d3a4330192c7fbfe93a8dbeb5"),
		DBAddress:  fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "5432")),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
