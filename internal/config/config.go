package config

import (
	"log"
	"os"
)

// Config содержит настройки приложения.
type Config struct {
	DatabaseURL string
	Port        string
}

// FromEnv читает переменные окружения и возвращает Config.
// Если переменная PORT не задана, используется значение по умолчанию 8080.
func FromEnv() Config {
	cfg := Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	log.Printf("config loaded: port=%s", cfg.Port)
	return cfg
}
