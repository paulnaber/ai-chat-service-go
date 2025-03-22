package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Config holds all configuration for the application
type Config struct {
	Environment string `envconfig:"ENVIRONMENT" default:"development"`
	Server      ServerConfig
	Database    DatabaseConfig
	CORS        CORSConfig
	Auth        AuthConfig
}

// ServerConfig holds all server-related configuration
type ServerConfig struct {
	Port string `envconfig:"SERVER_PORT" default:"3000"`
}

// DatabaseConfig holds all database-related configuration
type DatabaseConfig struct {
	Host         string `envconfig:"DB_HOST" default:"localhost"`
	Port         string `envconfig:"DB_PORT" default:"5432"`
	User         string `envconfig:"DB_USER" default:"postgres"`
	Password     string `envconfig:"DB_PASSWORD" default:"postgres"`
	Name         string `envconfig:"DB_NAME" default:"aichat"`
	SSLMode      string `envconfig:"DB_SSLMODE" default:"disable"`
	MigrationDir string `envconfig:"DB_MIGRATION_DIR" default:"migrations"`
	DatabaseUrl  string `envconfig:"DB_URL" default:"postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable"`
}

// GetConnectionString returns a PostgreSQL connection string
func (c DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.Name, c.SSLMode)
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
}

// AuthConfig holds authentication configuration for Keycloak
type AuthConfig struct {
	KeycloakURL  string `envconfig:"KEYCLOAK_URL" default:"http://localhost:8080"`
	Realm        string `envconfig:"KEYCLOAK_REALM" default:"ai-chat"`
	ClientID     string `envconfig:"KEYCLOAK_CLIENT_ID" default:"ai-chat-client"`
	ClientSecret string `envconfig:"KEYCLOAK_CLIENT_SECRET" default:""`
	PublicKey    string `envconfig:"KEYCLOAK_PUBLIC_KEY" default:""`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	godotenv.Load()

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
