package request

import "seabattle/internal/service/action"

type SetShip struct {
	Code  string       `json:"code"`
	Point action.Point `json:"point"`
}
