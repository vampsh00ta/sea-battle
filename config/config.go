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
		//RMQ  `yaml:"rabbitmq"`
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
	// RMQ -.
	RMQ struct {
		ServerExchange string `env-required:"true" yaml:"rpc_server_exchange" env:"RMQ_RPC_SERVER"`
		ClientExchange string `env-required:"true" yaml:"rpc_client_exchange" env:"RMQ_RPC_CLIENT"`
		URL            string `env-required:"true"                            env:"RMQ_URL"`
	}
)

//// NewConfig returns seabattle config.
//func New() (*Config, error) {
//	cfg := &Config{}
//	err := godotenv.Load(".env")
//	if err != nil {
//		return nil, err
//	}
//	currPath, err := os.Getwd()
//	if err != nil {
//		return nil, err
//	}
//	filePath := currPath + os.Getenv("path") + "/" + os.Getenv("env") + ".yml"
//	err = seabattleenv.ReadConfig(filePath, cfg)
//	if err != nil {
//		return nil, fmt.Errorf("config error: %w", err)
//	}
//
//	err = seabattleenv.ReadEnv(cfg)
//	if err != nil {
//		return nil, err
//	}
//
//	return cfg, nil
//}

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
