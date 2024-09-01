package config

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	conf *Config
	once sync.Once
)

type Config struct {
	PORT         string
	GIN_MODE     string
	N8N_BASE_URL string
	N8N_API_KEY  string
}

func Load() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Printf(".env file not found: %s", err.Error())
		}

		conf = &Config{
			PORT:         getEnv("PORT", "6000"),
			GIN_MODE:     getEnv("GIN_MODE", "debug"),
			N8N_BASE_URL: getEnv("N8N_BASE_URL", "http://localhost:5678"),
			N8N_API_KEY:  getEnv("N8N_API_KEY", ""),
		}
	})
	return conf
}

func Get() *Config {
	return conf
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
