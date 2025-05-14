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

// LoadConfig loads application configuration from a YAML file and environment variables.
// It searches for a file named "config.yaml" in the current and parent directories, allowing environment variables to override file values.
// Returns a pointer to the populated Config struct or an error if loading or unmarshaling fails.
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
