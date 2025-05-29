package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DB                  DatabaseConfig
	Server              ServerConfig
	GRPC                GRPCConfig
	GOOGLE_CLIENT_TOKEN string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	Port string
}

type GRPCConfig struct {
	Port string
}

var (
	cfg  *Config
	once sync.Once
)

// LoadConfig loads the configuration from environment variables
func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env") // optional, can be skipped in prod
	if err != nil {
		return nil, err
	}

	config := &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
		},
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
		},
		GRPC: GRPCConfig{
			Port: os.Getenv("GRPC_PORT"),
		},
		GOOGLE_CLIENT_TOKEN: os.Getenv("GOOGLE_CLIENT_TOKEN"),
	}

	return config, nil
}

// GetConfig returns a singleton config instance
func GetConfig() *Config {
	once.Do(func() {
		var err error
		cfg, err = LoadConfig()
		if err != nil {
			panic("Failed to load configuration: " + err.Error())
		}
	})
	return cfg
}
