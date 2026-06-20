package config

type Config struct {
	ServerPort string
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

var AppConfig Config

func InitConfig() {
	AppConfig = Config{
		ServerPort: "8080",
		Database: DatabaseConfig{
			Driver: "sqlite",
			DSN:    "jianli.db",
		},
	}
}
