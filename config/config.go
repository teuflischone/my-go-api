package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config ...
type Config struct {
	DatabaseURL    string
	CacheURL       string
	LoggerLevel    string
	ContextTimeout int
	JWTSecretKey   string
}

// LoadConfig will load config from environment variable
func LoadConfig() (config *Config) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	cacheURL := os.Getenv("CACHE_URL")
	loggerLevel := os.Getenv("LOGGER_LEVEL")
	contextTimeout, _ := strconv.Atoi(os.Getenv("CONTEXT_TIMEOUT"))
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	return &Config{
		DatabaseURL:    databaseURL,
		CacheURL:       cacheURL,
		LoggerLevel:    loggerLevel,
		ContextTimeout: contextTimeout,
		JWTSecretKey:   jwtSecretKey,
	}
}
