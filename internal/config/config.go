package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port string
	Mode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Attempt to read config file
	if err := viper.ReadInConfig(); err != nil {
		viper.AutomaticEnv()
	}

	// Parse JWT expiration duration
	jwtExpStr := viper.GetString("JWT_EXPIRATION")
	if jwtExpStr == "" {
		jwtExpStr = "24h"
	}
	jwtExpiration, err := time.ParseDuration(jwtExpStr)
	if err != nil {
		return nil, err
	}

	// Parse allowed origins
	allowedOriginsStr := viper.GetString("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		allowedOriginsStr = "http://localhost:5173"
	}
	allowedOrigins := strings.Split(allowedOriginsStr, ",")

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Mode: viper.GetString("GIN_MODE"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Expiration: jwtExpiration,
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
	}

	return config, nil
}
