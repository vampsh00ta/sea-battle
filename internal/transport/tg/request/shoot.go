package request

import (
	"seabattle/internal/service/action"
)

type Shoot struct {
	Code  string       `json:"code"`
	TgId  string       `json:"tgId"`
	Point action.Point `json:"point"`
}
