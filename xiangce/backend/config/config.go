package config

import (
	"os"
)

type Config struct {
	Port         string
	UploadDir    string
	ThumbnailDir string
	JWTSecret    string
	DBPath       string
}

func LoadConfig() *Config {
	return &Config{
		Port:         getEnv("PORT", "8080"),
		UploadDir:    getEnv("UPLOAD_DIR", "./uploads"),
		ThumbnailDir: getEnv("THUMBNAIL_DIR", "./uploads/thumbnails"),
		JWTSecret:    getEnv("JWT_SECRET", "xiangce-secret-key-2024"),
		DBPath:       getEnv("DB_PATH", "./xiangce.db"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
