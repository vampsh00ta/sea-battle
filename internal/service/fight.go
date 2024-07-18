package service

import (
	"context"
	"encoding/json"
	"errors"
	kafkago "github.com/segmentio/kafka-go"
	"math/rand"
	"seabattle/internal/entity"

	"seabattle/internal/service/rules"
)

//type Fight interface {
//	Shoot(ctx context.Context, req entity.Shoot) (entity.Fight, int, error)
//	//AddShip(ctx context.Context, tgId string, p1, p2 entity.Point) error
//	JoinFight(ctx context.Context, code, tgId string) (entity.Fight, error)
//	CreateFight(ctx context.Context, tgId string) (string, error)
//	SetShip(ctx context.Context, req entity.SetShip) (*entity.BattleField, int, error)
//	SetFieldQueryId(ctx context.Context, sessionId, tgId, queryId string, my bool) error
//	//SearchFight(ctx context.Context, tgId int) error
//	InitFightAction(ctx context.Context, token string) (*entity.Fight, error)
//}

func (s service) SearchFight(ctx context.Context, tgId int) error {
	//connection, err := rmq.OpenConnection("queue", "tcp", "localhost:6379", 1, nil)
	//if err != nil {
	//	return err
	//}
	//taskQueue, err := connection.OpenQueue("game_search")
	//if err != nil {
	//	return err
	//}
	//msg := entity.SearchFightMsg{Rating: rand.Intn(1000), TgID: rand.Intn(10000)}
	//taskBytes, err := json.Marshal(msg)
	//if err != nil {
	//	return err
	//}
	//if err := taskQueue.PublishBytes(taskBytes); err != nil {
	//	return err
	//}
	msg := entity.SearchFight{Rating: rand.Intn(1000), TgID: rand.Intn(10000)}
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := s.kafka.WriteMessages(ctx, kafkago.Message{Value: value}); err != nil {
		return err
	}
	return nil
}
func (s service) InitFightAction(ctx context.Context, token string) (*entity.Fight, error) {
	//sessionId, err := s.psql.GetSessionByCode(ctx, token)
	//if err != nil {
	//	return nil, err
	//}

	fight, err := s.mongo.GetFight(ctx, token)
	if err != nil {
		return nil, err
	}
	return &fight, nil
}

func (s service) SetFieldQueryId(ctx context.Context, sessionID, tgID, queryID string, my bool) error {
	return s.mongo.SetFieldQueryId(ctx, sessionID, tgID, queryID, my)
}

func (s service) SetShip(ctx context.Context, req entity.SetShip) (*entity.BattleField, int, error) {

	p, err := s.mongo.GetPoint(ctx, req.Code, req.TgId)
	if err != nil {
		return nil, -1, err
	}
	var res int

	b, err := s.mongo.GetBattleField(ctx, req.Code, req.TgId, true)
	if err != nil {
		return nil, -1, err
	}
	switch {
	case p.X == -1 && p.Y == -1:
		if err = s.mongo.SetPoint(ctx, req.Code, req.TgId, req.Point); err != nil {
			return nil, -1, err
		}

		res = rules.ShipSecondPoint
		return b, res, nil
	default:

		isReady, err := s.addShipEntity(b, p, req.Point)
		if err != nil {
			if err := s.mongo.SetPoint(ctx, req.Code, req.TgId, entity.Point{-1, -1}); err != nil {
				return nil, -1, err
			}
			res = rules.ShipFirstPoint
			return b, res, err
		}
		if err := s.mongo.SetPoint(ctx, req.Code, req.TgId, entity.Point{-1, -1}); err != nil {
			return nil, -1, err
		}
		if err := s.mongo.SetBattleField(ctx, req.Code, req.TgId, b, true); err != nil {
			return nil, -1, err
		}
		if err := s.setReady(ctx, &res, isReady, req.Code); err != nil {
			return nil, -1, err
		}
	}
	return b, res, nil

}
func (s service) CreateFight(ctx context.Context, tgId string) (string, error) {

	code := s.GetInviteCode()
	//if err := s.psql.AddSession(ctx, code, session); err != nil {
	//	return "", err
	//}
	if err := s.mongo.CreateFight(ctx, entity.Fight{
		Users: []entity.User{
			{
				TgId: tgId,
			},
		},
		SessionId: code,
	}); err != nil {
		return code, nil
	}

	return code, nil
}

func (s service) JoinFight(ctx context.Context, sessionID, tgId string) (entity.Fight, error) {

	preFight, err := s.mongo.GetFight(ctx, sessionID)
	if err != nil {
		return entity.Fight{}, err
	}
	if preFight.Stage == 1 {
		return entity.Fight{}, errors.New(rules.GameOnProgressErr)
	}
	if len(preFight.Users) == 2 {
		return entity.Fight{}, errors.New(rules.AlreadyJoined)
	}

	joinedUser := entity.User{
		TgId: tgId,
	}
	tgIds := []string{preFight.Users[0].TgId, joinedUser.TgId}
	firstTurn := tgIds[rand.Intn(2)]

	preFight.Turn = firstTurn
	preFight.Stage = rules.StagePick
	preFight.Users = append(preFight.Users, joinedUser)

	for i := range preFight.Users {
		s.setUserParams(&preFight.Users[i])

	}

	if err := s.mongo.UpdateFight(ctx, preFight); err != nil {
		return entity.Fight{}, err
	}
	return preFight, nil
}

func (s service) Shoot(ctx context.Context, req entity.Shoot) (entity.Fight, int, error) {

	fight, err := s.mongo.GetFight(ctx, req.Code)
	if err != nil {
		return entity.Fight{}, -1, err
	}
	switch {
	case fight.Stage == rules.StageEnd:
		return entity.Fight{}, -1, errors.New(rules.GameEndedErr)
	case fight.Turn != req.TgId:
		return entity.Fight{}, -1, errors.New(rules.NotYourTurnErr)
	}

	attacker, defender := getCurrRoles(&fight)

	res, err := s.shootEntity(attacker.EnemyField, attacker.MyField, req.Point.Y, req.Point.X)
	if err != nil {
		return entity.Fight{}, -1, err
	}
	switch res {
	case rules.Missed:
		fight.Turn = defender.TgId
	case rules.Lost:
		fight.State = rules.StageEnd

	}

	if err := s.mongo.UpdateFight(ctx, fight); err != nil {
		return entity.Fight{}, -1, err
	}

	return fight, res, nil

}
