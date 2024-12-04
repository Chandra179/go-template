package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Database struct {
	*sql.DB
}

type DatabaseConfig struct {
	Host       string
	User       string
	Password   string
	DBName     string
	SSLMode    string
	DriverName string
	Port       int
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *DatabaseConfig) (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	db, err := sql.Open(cfg.DriverName, connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Database{db}, nil
}

func (db *Database) Close() error {
	return db.DB.Close()
}
