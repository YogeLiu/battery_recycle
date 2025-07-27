package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

// Config holds the application configuration
type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Server   struct {
		Port string `yaml:"port"`
		Mode string `yaml:"mode"`
	} `yaml:"server"`
}

// LoadConfig loads configuration from a YAML file based on the environment
func LoadConfig() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "test" // default to test environment
	}

	configPath := fmt.Sprintf("config/config_%s.yaml", env)

	// Try to read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", configPath, err)
	}

	// Set default values if not provided in config file
	if config.Database.Host == "" {
		panic("DB_HOST is not set")
	}

	if config.Database.Port == "" {
		panic("DB_PORT is not set")
	}

	if config.Database.User == "" {
		panic("DB_USER is not set")
	}

	if config.Database.Password == "" {
		panic("DB_PASSWORD is not set")
	}

	if config.Database.Name == "" {
		panic("DB_NAME is not set")
	}

	if config.Server.Port == "" {
		panic("PORT is not set")
	}

	if config.Server.Mode == "" {
		panic("MODE is not set")
	}

	return &config, nil
}

// GetDSN returns the data source name for the database connection
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name)
}
