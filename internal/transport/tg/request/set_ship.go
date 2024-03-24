package request

import "seabattle/internal/service/entity"

type SetShip struct {
	Point entity.Point `json:"point"`

	Code string `json:"code"`
}
