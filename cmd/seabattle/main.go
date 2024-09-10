package main

import (
	"log"
	"seabattle/config"
	app "seabattle/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	cfg, err := config.New(configPath)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.NewPooling(cfg)

}
