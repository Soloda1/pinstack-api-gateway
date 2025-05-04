package config

import (
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
	Address string
	Port    int
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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("env", "dev")

	viper.SetDefault("http_server.address", "0.0.0.0")
	viper.SetDefault("http_server.port", 8080)

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

	config := &Config{
		Env: viper.GetString("env"),
		HTTPServer: HTTPServer{
			Address: viper.GetString("http_server.address"),
			Port:    viper.GetInt("http_server.port"),
		},
		Services: Services{
			User: UserService{
				Address: viper.GetString("services.user.address"),
				Port:    viper.GetInt("services.user.port"),
			},
			Auth: AuthService{
				Address: viper.GetString("services.auth.address"),
				Port:    viper.GetInt("services.auth.port"),
			},
		},
		JWT: JWT{
			Secret:           viper.GetString("jwt.secret"),
			AccessExpiresAt:  viper.GetString("jwt.access_expires_at"),
			RefreshExpiresAt: viper.GetString("jwt.refresh_expires_at"),
		},
	}

	return config
}
