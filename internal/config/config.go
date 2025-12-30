package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBPort      string
	DBSSLMODE   string
	CORSOrigin  string

	ServerPort  string
}

func Load() *Config {

	if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }

    c := &Config{
        DBHost:     getEnv("DB_HOST", ""),
        DBUser:     getEnv("DB_USER", ""),
        DBPassword: getEnv("DB_PASS", ""),
        DBName:     getEnv("DB_NAME", ""),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBSSLMODE:  getEnv("DB_SSLMODE", "require"),
		CORSOrigin: getEnv("CORS_ORIGIN", "http://127.0.0.1:5173"),


        ServerPort: getEnv("SERVER_PORT", "8080"),
    }

	if c.DBHost == "" || c.DBUser == "" || c.DBPassword == "" || c.DBName == "" {
		log.Fatal("Missing db ENV Variables!!!")
	}

	return c
}

func getEnv(key, defaultstr string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultstr
}