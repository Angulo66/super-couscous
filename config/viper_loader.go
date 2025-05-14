package config

import "github.com/spf13/viper"

type ViperLoader struct {
	configName string
	configType string
	paths      []string
	useEnv     bool
}

func NewViperLoader(configName, configType string, paths []string, useEnv bool) *ViperLoader {
	return &ViperLoader{
		configName: configName,
		configType: configType,
		paths:      paths,
		useEnv:     useEnv,
	}
}

func DefaultViperLoader() *ViperLoader {
	return NewViperLoader("config", "yaml", []string{".", ".."}, true)
}

func (v *ViperLoader) SetDefaults() error {
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.debug", false)
	viper.SetDefault("ratelimit.enabled", false)
	viper.SetDefault("ratelimit.max_requests", 100)
	viper.SetDefault("ratelimit.window_size", "1m")
	return nil
}

func (v *ViperLoader) Load(cfg interface{}) error {
	viper.SetConfigName(v.configName)
	viper.SetConfigType(v.configType)

	for _, path := range v.paths {
		viper.AddConfigPath(path)
	}

	if v.useEnv {
		viper.AutomaticEnv()
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return viper.Unmarshal(cfg)
}
