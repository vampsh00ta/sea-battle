package service

import (
	"context"
	"seabattle/internal/entity"
)

type Fight interface {
	Shoot(ctx context.Context, req entity.Shoot) (entity.Fight, int, error)
	//AddShip(ctx context.Context, tgId string, p1, p2 entity.Point) error
	JoinFight(ctx context.Context, code, tgId string) (entity.Fight, error)
	CreateFight(ctx context.Context, tgId string) (string, error)
	SetShip(ctx context.Context, req entity.SetShip) (*entity.BattleField, int, error)
	SetFieldQueryId(ctx context.Context, sessionId, tgId, queryId string, my bool) error
	//SearchFight(ctx context.Context, tgId int) error
	InitFightAction(ctx context.Context, token string) (*entity.Fight, error)
}
