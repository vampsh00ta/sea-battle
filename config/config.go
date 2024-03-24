package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	// Config -.
	Config struct {
		App   `yaml:"seabattle"`
		HTTP  `yaml:"http"`
		Log   `yaml:"logger"`
		PG    `yaml:"postgres"`
		Redis `yaml:"redis"`
		Tg    `yaml:"tg"`
		//RMQ  `yaml:"rabbitmq"`
	}
	Tg struct {
		Apitoken string `env-required:"true" yaml:"apitoken"    env:"API_TOKEN"`
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
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}

	// PG -.
	PG struct {
		PoolMax  int    `env-required:"true" yaml:"pool_max" env:"PG_POOL_MAX"`
		Username string `env-required:"true" yaml:"username" env-default:"postgres"`
		Password string `env-required:"true" yaml:"password" env-default:"postgres"`
		Host     string `env-required:"true" yaml:"host" env-default:"localhost"`
		Port     string `env-required:"true" yaml:"port" env-default:"5432"`
		Name     string `env-required:"true" yaml:"name" env-default:"postgres"`
	}
	// Redis -.
	Redis struct {
		Address  string `env-required:"true" yaml:"address" env:"address"`
		Password string `env-required:"true" yaml:"password" env-default:"password"`
		Db       int    `env-required:"true" yaml:"db" env-default:"db"`
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

func New() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	currPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := currPath + os.Getenv("path") + "/" + os.Getenv("env") + ".yml"
	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	var cfg *Config

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil

}

func NewGame() (*Game, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}
	currPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	filePath := currPath + os.Getenv("path") + "/" + "game.yaml"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	var cfg *Game

	if err := d.Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil

}
