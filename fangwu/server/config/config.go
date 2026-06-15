package config

import "os"

type AppConfig struct {
	JWTSecret string
	Port      string
	DBPath    string
}

var Cfg AppConfig

func Init() {
	Cfg = AppConfig{
		JWTSecret: getEnv("JWT_SECRET", "fangwu-rental-secret-key-2024"),
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "fangwu.db"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
