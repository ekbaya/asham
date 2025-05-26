package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB     DatabaseConfig
	Server ServerConfig
	GRPC   GRPCConfig
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

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
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
	}

	return config, nil
}
