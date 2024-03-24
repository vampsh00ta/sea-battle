package service

import (
	"context"
	"seabattle/internal/repository/models"
	"seabattle/internal/service/entity"
)

type Session interface {
	InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error)
}

func (s service) InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error) {
	var user1, user2 models.User
	user1.MyField = entity.NewBattleField()
	user1.EnemyField = entity.NewBattleField()
	user1.TgId = tg1

	user2.MyField = entity.NewBattleField()
	user2.EnemyField = entity.NewBattleField()
	user2.TgId = tg2

	session, err := s.redis.CreateSessionByChatId(ctx, tg1, tg2)
	if err != nil {
		return models.Fight{}, err
	}
	fightModel := models.Fight{User1: user1, User2: user2, Turn: tg1, SessionId: session}
	if err := s.redis.SetFight(ctx, fightModel); err != nil {
		return models.Fight{}, err
	}
	return fightModel, nil
}
