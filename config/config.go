package config

// Config holds all application configuration
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Auth      AuthConfig      `mapstructure:"auth"`
	RateLimit RateLimitConfig `mapstructure:"ratelimit"`
}

type ServerConfig struct {
	Port  int  `mapstructure:"port"`
	Debug bool `mapstructure:"debug"`
}

type DatabaseConfig struct {
	URL string `mapstructure:"url"` // Database URL
}

type AuthConfig struct {
	Token string `mapstructure:"token"` // Authentication token
}

type RateLimitConfig struct {
	Enabled     bool   `mapstructure:"enabled"`     // Enable rate limiter
	MaxRequests int    `mapstructure:"maxRequests"` // Maximum requests per interval
	WindowSize  string `mapstructure:"windowSize"`  // Interval for rate limit (e.g., 1m)
}

func LoadConfig(loader ConfigLoader) (*Config, error) {
	if loader == nil {
		loader = DefaultViperLoader()
	}

	loader.SetDefaults()

	var config Config
	if err := loader.Load(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadDefaultConfig() (*Config, error) {
	return LoadConfig(nil)
}
