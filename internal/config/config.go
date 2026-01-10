package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds application configuration
type Config struct {
	SupercellAPIKey  string
	APIToken         string
	PostgresHost     string
	PostgresPort     int
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
	TopPlayersLimit  int
	APIPort          int
	RetentionDays    int
	LogLevel         string
}

// LoadFromEnv loads configuration from environment variables
func LoadFromEnv() (*Config, error) {
	cfg := &Config{
		SupercellAPIKey:  getEnv("SUPERCELL_API_KEY", ""),
		APIToken:         getEnv("API_TOKEN", ""),
		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnvInt("POSTGRES_PORT", 5432),
		PostgresDB:       getEnv("POSTGRES_DB", "royale_api"),
		PostgresUser:     getEnv("POSTGRES_USER", "royale"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", ""),
		TopPlayersLimit:  getEnvInt("TOP_PLAYERS_LIMIT", 1000),
		APIPort:          getEnvInt("API_PORT", 8080),
		RetentionDays:    getEnvInt("RETENTION_DAYS", 7),
		LogLevel:         getEnv("LOG_LEVEL", "info"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate checks if required configuration values are set
func (c *Config) Validate() error {
	if c.SupercellAPIKey == "" {
		return fmt.Errorf("SUPERCELL_API_KEY is required")
	}
	if c.APIToken == "" {
		return fmt.Errorf("API_TOKEN is required")
	}
	if c.PostgresPassword == "" {
		return fmt.Errorf("POSTGRES_PASSWORD is required")
	}
	if c.TopPlayersLimit < 1 || c.TopPlayersLimit > 1000 {
		return fmt.Errorf("TOP_PLAYERS_LIMIT must be between 1 and 1000")
	}
	return nil
}

// PostgresDSN returns the PostgreSQL connection string
func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
