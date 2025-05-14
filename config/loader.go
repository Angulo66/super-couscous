package config

type ConfigLoader interface {
	Load(cfg interface{}) error
	SetDefaults() error
}
