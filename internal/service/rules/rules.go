package rules

import "seabattle/config"

const configPath = "config/game.yaml"

var Game, _ = config.NewGame(configPath)
