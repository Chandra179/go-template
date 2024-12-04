package configs

import (
	"os"
	"strconv"

	"github.com/Chandra179/go-template/internal/db"
	_ "github.com/lib/pq"
)

type Config struct {
	DbConfig *db.DatabaseConfig
}

func LoadConfigFromEnv() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432
	}

	return &Config{
		DbConfig: &db.DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     port,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
		},
	}, nil
}
