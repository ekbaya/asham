package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DB                   DatabaseConfig
	Server               ServerConfig
	GRPC                 GRPCConfig
	GOOGLE_CLIENT_TOKEN  string
	AZURE_TENANT_ID      string
	AZURE_CLIENT_ID      string
	AZURE_CLIENT_SECRET  string
	AZURE_USER_EMAIL     string
	EmailConfig          EmailConfig
	DOC_TEMPLATE_PATH    string
	ONEDRIVE_FOLDER_NAME string
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

type EmailConfig struct {
	Host              string
	Port              string
	Username          string
	Password          string
	From              string
	EmailTemplatePath string
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
		AZURE_TENANT_ID:     os.Getenv("AZURE_TENANT_ID"),
		AZURE_CLIENT_ID:     os.Getenv("AZURE_CLIENT_ID"),
		AZURE_CLIENT_SECRET: os.Getenv("AZURE_CLIENT_SECRET"),
		AZURE_USER_EMAIL:    os.Getenv("AZURE_USER_EMAIL"),
		EmailConfig: EmailConfig{
			Host:              os.Getenv("EMAIL_HOST"),
			Port:              os.Getenv("EMAIL_PORT"),
			Username:          os.Getenv("EMAIL_USERNAME"),
			Password:          os.Getenv("EMAIL_PASSWORD"),
			From:              os.Getenv("EMAIL_FROM"),
			EmailTemplatePath: "templates/welcome_email.html",
		},
		DOC_TEMPLATE_PATH:    "../templates/project_template.docx",
		ONEDRIVE_FOLDER_NAME: "ASHAM_ARSO_PLATFORM",
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
