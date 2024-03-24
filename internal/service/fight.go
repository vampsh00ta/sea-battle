package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"seabattle/internal/repository/models"
	"seabattle/internal/service/action"
	//"seabattle/internal/service/action"
	"seabattle/internal/service/rules"
	"seabattle/internal/transport/tg/request"
)

type Fight interface {
	Shoot(ctx context.Context, req request.Shoot) (models.Fight, int, error)
	//AddShip(ctx context.Context, tgId string, p1, p2 entity.Point) error
	JoinFight(ctx context.Context, code, tgId string) (models.Fight, error)
	CreateFight(ctx context.Context, tgId string) (string, error)
	SetShip(ctx context.Context, tgId string, point action.Point, token string) (*models.BattleField, int, error)
	SetFieldQueryId(ctx context.Context, tgId, queryId string, my bool) error
	//EmptyField() *models.BattleField

	InitFightAction(ctx context.Context, token string) (*models.Fight, error)
}

func (s service) InitFightAction(ctx context.Context, token string) (*models.Fight, error) {
	sessionId, err := s.psql.GetSessionByCode(ctx, token)
	if err != nil {
		return nil, err
	}

	fight, err := s.redis.GetFight(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return &fight, nil
}

func (s service) SetFieldQueryId(ctx context.Context, tgId, queryId string, my bool) error {
	return s.redis.SetFieldQueryId(ctx, tgId, queryId, my)
}

func (s service) SetShip(ctx context.Context, tgId string, point action.Point, token string) (*models.BattleField, int, error) {
	x0, y0, err := s.redis.GetPoint(ctx, tgId)
	var res int
	if err != nil {
		return nil, -1, err
	}

	po := action.Point{X: x0, Y: y0}
	b, err := s.redis.GetBattleField(ctx, tgId, true)
	if err != nil {
		return nil, -1, err
	}
	if x0 == -1 && y0 == -1 {
		if err := s.redis.SetPoint(ctx, tgId, point.X, point.Y); err != nil {
			return nil, -1, err
		}
		res = rules.ShipSecondPoint
		return b, res, nil
	} else {
		isReady, err := s.action.AddShip(b, po, point)
		if err != nil {
			if err := s.redis.SetPoint(ctx, tgId, -1, -1); err != nil {
				return nil, -1, err
			}
			res = rules.ShipFirstPoint
			fmt.Println(err, "set ship")
			return b, res, err
		}
		if err := s.redis.SetPoint(ctx, tgId, -1, -1); err != nil {
			return nil, -1, err
		}
		if err := s.redis.SetBattleField(ctx, tgId, b, true); err != nil {
			return nil, -1, err
		}
		if err := s.setReady(ctx, &res, isReady, token); err != nil {
			return nil, -1, err
		}
	}
	return b, res, nil

}
func (s service) CreateFight(ctx context.Context, tgId string) (string, error) {
	session, err := s.redis.CreateSessionOnePerson(ctx, tgId)
	if err != nil {
		return "", err
	}
	code := s.GetInviteCode()
	if err := s.psql.AddSession(ctx, code, session); err != nil {
		return "", err
	}

	return code, nil
}

func (s service) JoinFight(ctx context.Context, code, tgId string) (models.Fight, error) {
	sessionId, err := s.psql.GetSessionByCode(ctx, code)
	if err != nil {
		return models.Fight{}, err
	}
	session, err := s.redis.GetSession(ctx, sessionId)
	if err != nil {
		return models.Fight{}, err
	}
	if session.Stage == 1 {
		return models.Fight{}, errors.New(rules.GameOnProgressErr)
	}
	if session.TgId1 == tgId {
		return models.Fight{}, errors.New(rules.AlreadyJoined)
	}

	session.TgId2 = tgId
	tgIds := []string{session.TgId1, session.TgId2}
	firstTurn := tgIds[rand.Intn(2)]
	session.Turn = firstTurn
	session.Ready = 0
	session.Stage = rules.StagePick
	if err != nil {
		return models.Fight{}, err
	}
	if err := s.redis.SetSession(ctx, sessionId, session); err != nil {
		return models.Fight{}, err
	}

	var user1, user2 models.User
	s.setUserParams(session.TgId1, &user1)
	s.setUserParams(session.TgId2, &user2)
	fightModel := models.Fight{
		User1:     user1,
		User2:     user2,
		Turn:      session.Turn,
		SessionId: sessionId,
		State:     rules.StagePick}
	if err := s.redis.SetFight(ctx, fightModel); err != nil {
		return models.Fight{}, err
	}
	return fightModel, nil
}

//	func (s service) AddShip(ctx context.Context, tgId string, p1, p2 entity.Point) error {
//		b, err := s.redis.GetBattleField(ctx, tgId, true)
//		if err != nil {
//			return err
//		}
//		if err := entity.AddShip(b, p1, p2); err != nil {
//			return err
//		}
//
//		if err := s.redis.SetBattleField(ctx, tgId, b, true); err != nil {
//			return err
//		}
//		return nil
//
// }
func (s service) Shoot(ctx context.Context, req request.Shoot) (models.Fight, int, error) {
	sessionId, err := s.psql.GetSessionByCode(ctx, req.Code)
	if err != nil {
		return models.Fight{}, -1, err
	}
	session, err := s.redis.GetSession(ctx, sessionId)
	if err != nil {
		return models.Fight{}, -1, err
	}
	if session.Stage == rules.StageEnd {
		return models.Fight{}, -1, errors.New(rules.GameEndedErr)
	}
	if session.Turn != req.TgId {
		return models.Fight{}, -1, errors.New(rules.NotYourTurnErr)
	}
	fight, err := s.redis.GetFight(ctx, sessionId)
	getAttacker(&fight)

	res, err := s.action.Shoot(fight.User1.EnemyField, fight.User2.MyField, req.Point.Y, req.Point.X)
	if err != nil {
		return models.Fight{}, -1, err
	}
	switch res {
	case rules.Missed:

		//turn := fight.User2.TgId
		fight.Turn = fight.User2.TgId
		session.Turn = fight.User2.TgId
	case rules.Lost:
		fight.State = rules.StageEnd
		session.Stage = rules.StageEnd

	}

	if err := s.redis.SetSession(ctx, sessionId, session); err != nil {
		return models.Fight{}, -1, err
	}
	if err := s.redis.SetFight(ctx, fight); err != nil {
		return models.Fight{}, -1, err
	}

	return fight, res, nil

}
