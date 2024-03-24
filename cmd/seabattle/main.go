package main

import (
	"log"
	"seabattle/config"
	app "seabattle/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	//r, err := rand.Int(rand.Reader, big.NewInt(2))
	//fmt.Println(r, err)
	app.NewPooling(cfg)

}
