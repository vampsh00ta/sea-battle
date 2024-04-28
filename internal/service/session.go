package service

import (
	"context"
	"seabattle/internal/models"
)

type Session interface {
	InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error)
}

func (s service) InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error) {
	var user1, user2 models.User
	user1.MyField = s.action.NewBattleField()
	user1.EnemyField = s.action.NewBattleField()
	user1.TgId = tg1

	user2.MyField = s.action.NewBattleField()
	user2.EnemyField = s.action.NewBattleField()
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
