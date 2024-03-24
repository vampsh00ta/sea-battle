package action

import (
	"seabattle/config"
)

type action struct {
	cfg *config.Game
}

type Action interface {
	Field
}

func New(cfg *config.Game) Action {
	return action{
		cfg: cfg,
	}
}
