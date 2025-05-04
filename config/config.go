package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string
	HTTPServer HTTPServer
	Services   Services
	JWT        JWT
}

type HTTPServer struct {
	Port        int    `mapstructure:"port"`
	Address     string `mapstructure:"address"`
	Timeout     int    `mapstructure:"timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
}

type Services struct {
	User UserService
	Auth AuthService
}

type UserService struct {
	Address string
	Port    int
}

type AuthService struct {
	Address string
	Port    int
}

type JWT struct {
	Secret           string
	AccessExpiresAt  string
	RefreshExpiresAt string
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config"
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)

	viper.SetDefault("env", "dev")

	viper.SetDefault("http_server.address", "0.0.0.0")
	viper.SetDefault("http_server.port", 8080)
	viper.SetDefault("http_server.timeout", 60)
	viper.SetDefault("http_server.idle_timeout", 120)

	viper.SetDefault("services.user.address", "user-service")
	viper.SetDefault("services.user.port", 50051)

	viper.SetDefault("services.auth.address", "auth-service")
	viper.SetDefault("services.auth.port", 50052)

	viper.SetDefault("jwt.secret", "my-secret")
	viper.SetDefault("jwt.access_expires_at", "1m")
	viper.SetDefault("jwt.refresh_expires_at", "5m")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file: %s", err)
		os.Exit(1)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(fmt.Errorf("error unmarshaling config: %w", err))
	}

	return &cfg
}
