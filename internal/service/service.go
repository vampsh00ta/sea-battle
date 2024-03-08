package service

import (
	"context"
	"errors"
	"seabattle/internal/repository/models"
	"seabattle/internal/repository/redis"
	"seabattle/internal/service/entity"
)

type Service interface {
	Fight
	Session
}
type Fight interface {
	Shoot(ctx context.Context, fight entity.Fight, p entity.Point) (models.Fight, int, error)
	AddShip(ctx context.Context, tgId string, p1, p2 entity.Point, shipType int) error
}

type Session interface {
	InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error)
}
type service struct {
	repo redis.Repository
}

func (s service) AddShip(ctx context.Context, tgId string, p1, p2 entity.Point, shipType int) error {
	b, err := s.repo.GetBattleField(ctx, tgId, true)
	if err != nil {
		return err
	}
	if err := entity.AddShip(b, p1, p2, shipType); err != nil {
		return err
	}

	if err := s.repo.SetBattleField(ctx, tgId, b, true); err != nil {
		return err
	}
	return nil

}
func (s service) Shoot(ctx context.Context, fight entity.Fight, p entity.Point) (models.Fight, int, error) {
	session, err := s.repo.GetSession(ctx, fight.SessionId)
	if err != nil {
		return models.Fight{}, -1, err
	}
	if session.Stage == entity.StageEnd {
		return models.Fight{}, -1, errors.New(entity.GameEndedErr)
	}
	if session.Turn != fight.Attacker {
		return models.Fight{}, -1, errors.New(entity.NotYourTurnErr)
	}

	fightModel, err := s.repo.GetFight(ctx, fight.SessionId)
	res, err := entity.Shoot(fightModel.Attacker.EnemyField, fightModel.Defender.MyField, p.Y, p.X)
	if err != nil {
		return models.Fight{}, -1, err
	}
	switch res {
	case entity.Missed:
		fight.Turn = fight.Defender
		session.Turn = fight.Defender
	case entity.Lost:
		fight.Stage = entity.StageEnd
		session.Stage = entity.StageEnd

	}

	if err := s.repo.SetSession(ctx, fight.SessionId, session); err != nil {
		return models.Fight{}, -1, err
	}
	if err := s.repo.SetFight(ctx, fightModel); err != nil {
		return models.Fight{}, -1, err
	}

	return fightModel, res, nil

}

func (s service) InitSessionFight(ctx context.Context, tg1, tg2 string) (models.Fight, error) {
	var attacker, defender models.User
	attacker.MyField = entity.NewBattleField()
	attacker.EnemyField = entity.NewBattleField()
	attacker.TgId = tg1

	defender.MyField = entity.NewBattleField()
	defender.EnemyField = entity.NewBattleField()
	defender.TgId = tg2

	session, err := s.repo.CreateSessionByChatId(ctx, tg1, tg2)
	if err != nil {
		return models.Fight{}, err
	}
	fightModel := models.Fight{Defender: defender, Attacker: attacker, Turn: tg1, SessionId: session}
	if err := s.repo.SetFight(ctx, fightModel); err != nil {
		return models.Fight{}, err
	}
	return fightModel, nil
}
func New(repo redis.Repository) Service {
	return &service{
		repo: repo,
	}
}
