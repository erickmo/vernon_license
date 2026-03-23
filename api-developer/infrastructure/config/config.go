package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	HTTPPort    string `mapstructure:"HTTP_PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
	JWTExpHours int    `mapstructure:"JWT_EXP_HOURS"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	viper.SetDefault("HTTP_PORT", "8081")
	viper.SetDefault("JWT_EXP_HOURS", 8)
	viper.SetDefault("LOG_LEVEL", "info")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("peringatan: .env tidak ditemukan, menggunakan env var")
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("gagal load config: %v", err)
	}
	return &cfg
}
