package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

type Config struct {
	Database DatabaseConfig
	JWT      JWTConfig
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")     // local/dev
	viper.AddConfigPath("/app")  // Docker WORKDIR
	viper.AddConfigPath("../..") // project root (for tests)
	viper.AutomaticEnv()

	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("jwt.secret", "JWT_SECRET_KEY")
	viper.BindEnv("jwt.expiration", "JWT_EXPIRE_HOURS")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &cfg, nil
}
