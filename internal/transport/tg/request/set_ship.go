package request

import (
	"seabattle/internal/entity"
)

type SetShip struct {
	Code  string       `json:"code"`
	Point entity.Point `json:"point"`
}
