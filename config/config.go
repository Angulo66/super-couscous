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

// LoadConfig loads application configuration using the provided ConfigLoader.
// If loader is nil, it uses the DefaultViperLoader.
//
// Example usage:
//
//	// Use the default config loader
//	cfg, err := config.LoadDefaultConfig()
//	if err != nil {
//		log.Fatal().Msgf("failed to load config %v", err)
//	}
//
//	// Or use a custom loader
//	loader := config.NewViperLoader("custom-config", "json", []string{"./config"}, true)
//	cfg, err := config.LoadConfig(loader)
//	if err != nil {
//		log.Fatal().Msgf("failed to load config %v", err)
//	}
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
