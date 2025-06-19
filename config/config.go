package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Env        string     `mapstructure:"env"`
	HTTPServer HTTPServer `mapstructure:"http_server"`
	Services   Services   `mapstructure:"services"`
	JWT        JWT        `mapstructure:"jwt"`
}

type HTTPServer struct {
	Port        int    `mapstructure:"port"`
	Address     string `mapstructure:"address"`
	Timeout     int    `mapstructure:"timeout"`
	IdleTimeout int    `mapstructure:"idle_timeout"`
}

type Services struct {
	User     UserService     `mapstructure:"user"`
	Auth     AuthService     `mapstructure:"auth"`
	Post     PostService     `mapstructure:"post"`
	Relation RelationService `mapstructure:"relation"`
}

type UserService struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type AuthService struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type PostService struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type RelationService struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

type JWT struct {
	Secret           string `mapstructure:"secret"`
	AccessExpiresAt  string `mapstructure:"access_expires_at"`
	RefreshExpiresAt string `mapstructure:"refresh_expires_at"`
}

func MustLoad() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	viper.SetDefault("env", "dev")

	viper.SetDefault("http_server.address", "0.0.0.0")
	viper.SetDefault("http_server.port", 8080)
	viper.SetDefault("http_server.timeout", 60)
	viper.SetDefault("http_server.idle_timeout", 120)

	viper.SetDefault("services.user.address", "user-service")
	viper.SetDefault("services.user.port", 50051)

	viper.SetDefault("services.auth.address", "auth-service")
	viper.SetDefault("services.auth.port", 50052)

	viper.SetDefault("services.auth.address", "post-service")
	viper.SetDefault("services.auth.port", 50053)

	viper.SetDefault("services.relation.address", "relation-service")
	viper.SetDefault("services.relation.port", 50054)

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
