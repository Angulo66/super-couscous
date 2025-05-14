package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port  int
		Debug bool
	}
	Database struct {
		URL string
	}
	Auth struct {
		Token string
	}
	RateLimit struct {
		Enabled     bool
		MaxRequests int
		WindowSize  string
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv() // Override with env vars if set

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
