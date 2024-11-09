package configs

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Database struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// LoadConfigFromEnv loads database configuration from environment variables
func LoadConfigFromEnv() (*Config, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		port = 5432 // default port
	}

	return &Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     port,
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
	}, nil
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *Config) (*Database, error) {
	// Construct connection string
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Create context with timeout for connection test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Verify connection
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	return &Database{db}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}
