package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	RedisURL           string
	Port               string
	WorkerCount        int
	MaxRetries         int
	RetryBaseDelayS    int
	RateLimitRPM       int
	BankFailureRate    float64
	HighValueThreshold int64
	MLServiceURL       string
	MLTimeoutMS        int
	VelocityWindowS    int
	VelocityMaxCount   int
}

func mustGet(key string) string {
	val := os.Getenv(key)

	if val == "" {
		log.Fatalf("FATAL: required environment variable %s is not set", key)
	}
	return val
}

func getOrDefault(key string, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getIntOrDefault(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		n, err := strconv.Atoi(val)
		if err != nil {
			log.Fatalf("FATAL: %s must be an integer, got: %s", key, val)
		}
		return n
	}
	return defaultVal

}

func getFloatOrDefault(key string, defaultVal float64) float64 {
	if val := os.Getenv(key); val != "" {
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			log.Fatalf("FATAL: %s must be an float, got: %s", key, val)
		}
		return f
	}
	return defaultVal
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		DatabaseURL:        mustGet("DATABASE_URL"),
		RedisURL:           mustGet("REDIS_URL"),
		Port:               getOrDefault("PORT", "8080"),
		WorkerCount:        getIntOrDefault("WORKER_COUNT", 5),
		MaxRetries:         getIntOrDefault("MAX_RETRIES", 5),
		RetryBaseDelayS:    getIntOrDefault("RETRY_BASE_DELAY_S", 30),
		RateLimitRPM:       getIntOrDefault("RATE_LIMIT_RPM", 100),
		BankFailureRate:    getFloatOrDefault("BANK_FAILURE_RATE", 0.2),
		HighValueThreshold: int64(getIntOrDefault("HIGH_VALUE_THRESHOLD", 10000000)),
		MLServiceURL:       os.Getenv("ML_SERVICE_URL"),
		MLTimeoutMS:        getIntOrDefault("ML_TIMEOUT_MS", 100),
		VelocityWindowS:    getIntOrDefault("VELOCITY_WINDOW_S", 60),
		VelocityMaxCount:   getIntOrDefault("VELOCITY_MAX_COUNT", 5),
	}
	return cfg
}
