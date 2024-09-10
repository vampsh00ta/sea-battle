package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"path"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"app"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		Mongo `yaml:"mongo"`
		Tg    `yaml:"tg"`
	}
	SearchService struct {
		Host string `env-required:"true" yaml:"port" env:"SEARCH_SERVICE_HOST"`
		Port string `env-required:"true" yaml:"port" env:"SEARCH_SERVICE_PORT"`
	}
	Tg struct {
		ApiToken string `env-required:"true" yaml:"apitoken"    env:"API_TOKEN"`
		BaseURL  string `env-required:"true" yaml:"baseURL"    env:"API_TOKEN"`
	}
	// App -.
	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	// HTTP -.
	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	}

	// Log -.
	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// Mongo -.

	Mongo struct {
		Address    string `env-required:"true" yaml:"address" env:"MONGO_HOST"`
		Db         string `env-required:"true" yaml:"db" env:"MONGO_DB"`
		Collection string `env-required:"true" yaml:"collection" env:"MONGO_COLLECTION"`
		URL        string `env-required:"true"                      env:"MONGO_URL"`
	}
)
type (
	Game struct {
		Main `yaml:"main"`
	}

	Main struct {
		Height        int `env-required:"true" yaml:"height" env:"height"`
		Weight        int `env-required:"true" yaml:"weight" env:"weight"`
		MaxShipCount  int `env-required:"true" yaml:"max_ship_count" env:"max_ship_count"`
		ShipTypeCount int `env-required:"true" yaml:"ship_type_count" env:"ship_type_count"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}

func NewGame(configPath string) (*Game, error) {

	cfg := &Game{}

	err := cleanenv.ReadConfig(path.Join("./", configPath), cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil

}
