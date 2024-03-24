package request

import (
	"seabattle/internal/service/action"
)

type SetShip struct {
	Point action.Point `json:"point"`

	Code string `json:"code"`
}
