package config

import "os"

type Config struct {
	Port    string
	DBPath  string
	DataDir string
}

func Default() Config {
	return Config{
		Port:    getenv("QRFORGE_PORT", "8080"),
		DBPath:  getenv("QRFORGE_DB", "./data/qrcode.db"),
		DataDir: getenv("QRFORGE_DATA", "./data"),
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
