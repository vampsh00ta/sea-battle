package main

import (
	"log"
	"seabattle/config"
	app "seabattle/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.NewPooling(cfg)

}
