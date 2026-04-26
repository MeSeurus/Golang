package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	ServerPort         string
	JWTSecret          string
	JWTExpirationHours int

	SchedulerInterval time.Duration
	WorkerPoolSize    int
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	expHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if err != nil {
		expHours = 24
	}

	workerPoolSize, err := strconv.Atoi(getEnv("WORKER_POOL_SIZE", "5"))
	if err != nil {
		workerPoolSize = 5
	}

	schedInterval, err := time.ParseDuration(getEnv("SCHEDULER_INTERVAL", "10s"))
	if err != nil {
		schedInterval = 10 * time.Second
	}

	return &Config{
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5435"),
		DBUser:             getEnv("DB_USER", "blog_user"),
		DBPassword:         getEnv("DB_PASSWORD", "blog_pass"),
		DBName:             getEnv("DB_NAME", "blog"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		ServerPort:         getEnv("SERVER_PORT", "8083"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me"),
		JWTExpirationHours: expHours,
		SchedulerInterval:  schedInterval,
		WorkerPoolSize:     workerPoolSize,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
