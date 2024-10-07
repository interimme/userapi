package config

import (
	"log"
	"os"
	"strconv"
)

// Config structure holds all configuration values for the application.
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig holds the database-related configuration.
type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     int
}

// ServerConfig holds the server-related configuration.
type ServerConfig struct {
	HttpPort int
	GrpcPort int
	GinPort  int
}

// Init initializes the configuration by reading from environment variables.
func Init() *Config {
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("DB_PORT doesn't look like an integer: %s", err)
	}
	httpPort, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatalf("HTTP_PORT doesn't look like an integer: %s", err)
	}
	grpcPort, err := strconv.Atoi(os.Getenv("GRPC_PORT"))
	if err != nil {
		log.Fatalf("GRPC_PORT doesn't look like an integer: %s", err)
	}
	ginPort, err := strconv.Atoi(os.Getenv("GIN_PORT"))
	if err != nil {
		log.Fatalf("GIN_PORT doesn't look like an integer: %s", err)
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			Port:     dbPort,
		},
		Server: ServerConfig{
			HttpPort: httpPort,
			GrpcPort: grpcPort,
			GinPort:  ginPort,
		},
	}
}
