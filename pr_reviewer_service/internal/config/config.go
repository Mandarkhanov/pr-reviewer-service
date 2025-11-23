package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServerAddress string
	DBUser        string
	DBPassword    string
	DBName        string
	DBHost        string
	DBPort        string
}

func Load() (*Config, error) {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	httpPort := os.Getenv("APP_HTTP_PORT")

	if dbUser == "" || dbPassword == "" || dbName == "" || dbHost == "" || dbPort == "" || httpPort == "" {
		return nil, fmt.Errorf("one or more environment variables are missing")
	}

	cfg := &Config{
		ServerAddress: ":" + httpPort,
		DBUser:        dbUser,
		DBPassword:    dbPassword,
		DBName:        dbName,
		DBHost:        dbHost,
		DBPort:        dbPort,
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}
