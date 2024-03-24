package request

import "seabattle/internal/service/entity"

type Shoot struct {
	Code  string       `json:"code"`
	TgId  string       `json:"tgId"`
	Point entity.Point `json:"point"`
}
