package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	// Server configuration
	Server struct {
		Port  int  `mapstructure:"port"`
		Debug bool `mapstructure:"debug"`
	}
	// Database configuration
	Database struct {
		URL string `mapstructure:"url"` // Database URL
	}
	// Authentication configuration
	Auth struct {
		Token string `mapstructure:"token"` // Authentication token
	}
	// Rate limit configuration
	RateLimit struct {
		Enabled     bool   `mapstructure:"enabled"`     // Enable rate limiter
		MaxRequests int    `mapstructure:"maxRequests"` // Maximum requests per interval
		WindowSize  string `mapstructure:"windowSize"`  // Interval for rate limit (e.g., 1m)
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv() // Override with env vars if set

	// Set default values
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.debug", false)
	viper.SetDefault("ratelimit.enabled", false)
	viper.SetDefault("ratelimit.max_requests", 100)
	viper.SetDefault("ratelimit.window_size", "1m")

	var config Config
	if err := viper.ReadInConfig(); err != nil {
		// Allow missing config file in developement
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
		log.Warn().Msg("No configuration file found, using defaults and environment variables")
	}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
