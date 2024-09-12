package service

import (
	"context"
	"errors"
	"math/rand"
	"seabattle/config"
	irep "seabattle/internal/app/repository/mongorep"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/entity"

	"seabattle/internal/service/rules"
)

//	type Fight interface {
//		Shoot(ctx context.Context, req entity.Shoot) (entity.Fight, int, error)
//		//AddShip(ctx context.Context, tgId string, p1, p2 entity.Point) error
//		JoinFight(ctx context.Context, code, tgId string) (entity.Fight, error)
//		CreateFight(ctx context.Context, tgId string) (string, error)
//		SetShip(ctx context.Context, req entity.SetShip) (*entity.BattleField, int, error)
//		SetFieldQueryId(ctx context.Context, sessionId, tgId, queryId string, my bool) error
//		//SearchFight(ctx context.Context, tgId int) error
//		InitFightAction(ctx context.Context, token string) (*entity.Fight, error)
//	}

type battlePreparation struct {
	fight         irep.Fight
	user          irep.User
	battleField   irep.BattleField
	entity        EntityInteraction
	helpers       Helpers
	codeGenerator CodeGenerator
}

func NewBattlePreparation(
	fightRepo irep.Fight,
	userRepo irep.User,
	battleField irep.BattleField,
	game *config.Game) isrvc.BattlePreparation {

	return &battlePreparation{
		fight:         fightRepo,
		user:          userRepo,
		battleField:   battleField,
		entity:        newEntityInteraction(game),
		helpers:       newHelpers(game),
		codeGenerator: newCodeGenerator(),
	}
}

func (s battlePreparation) setReady(ctx context.Context, res *int, ready int, token string) error {
	if ready != rules.PersonsReady {
		return nil
	}
	fight, err := s.fight.GetBySessionID(ctx, token)

	if err != nil {
		return err
	}
	if ready == rules.PersonsReady {
		fight.Ready += 1

	}
	switch fight.Ready {
	case 1:
		*res = rules.PersonReady
	case 2:
		*res = rules.PersonsReady

	}
	if err := s.fight.Update(ctx, fight); err != nil {

		return err
	}

	return nil
}
func (s battlePreparation) InitFightAction(ctx context.Context, token string) (*entity.Fight, error) {
	//sessionId, err := s.psql.GetSessionByCode(ctx, token)
	//if err != nil {
	//	return nil, err
	//}

	f, err := s.fight.GetBySessionID(ctx, token)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (s battlePreparation) SetFieldQueryId(ctx context.Context, sessionID, tgID, queryID string, my bool) error {
	return s.user.SetFieldQueryId(ctx, sessionID, tgID, queryID, my)
}

func (s battlePreparation) SetShip(ctx context.Context, req entity.SetShip) (*entity.BattleField, int, error) {

	p, err := s.user.GetPoint(ctx, req.Code, req.TgId)
	if err != nil {
		return nil, -1, err
	}
	var res int

	b, err := s.battleField.GetBySessionID(ctx, req.Code, req.TgId, true)
	if err != nil {
		return nil, -1, err
	}
	switch {
	case p.X == -1 && p.Y == -1:
		if err = s.user.SetPoint(ctx, req.Code, req.TgId, req.Point); err != nil {
			return nil, -1, err
		}

		res = rules.ShipSecondPoint
		return b, res, nil
	default:

		isReady, err := s.entity.addShip(b, p, req.Point)
		if err != nil {
			if setPointErr := s.user.SetPoint(ctx, req.Code, req.TgId, entity.Point{X: -1, Y: -1}); setPointErr != nil {
				return nil, -1, setPointErr
			}
			res = rules.ShipFirstPoint
			return b, res, err
		}
		if err := s.user.SetPoint(ctx, req.Code, req.TgId, entity.Point{X: -1, Y: -1}); err != nil {
			return nil, -1, err
		}
		if err := s.battleField.Set(ctx, req.Code, req.TgId, b, true); err != nil {
			return nil, -1, err
		}
		if err := s.setReady(ctx, &res, isReady, req.Code); err != nil {
			return nil, -1, err
		}
	}
	return b, res, nil

}

func (s battlePreparation) InitFight(ctx context.Context, tgIDs ...string) (entity.Fight, error) {

	code := s.codeGenerator.GetCode()

	var fight entity.Fight

	firstTurn := tgIDs[rand.Intn(2)]

	fight.Turn = firstTurn
	fight.Stage = rules.StagePick
	fight.SessionId = code

	for _, tgID := range tgIDs {
		var user entity.User
		user.TgId = tgID
		s.helpers.setUserParams(&user)
		fight.Users = append(fight.Users, user)
	}

	if err := s.fight.Create(ctx, fight); err != nil {
		return entity.Fight{}, err
	}
	return fight, nil
}

func (s battlePreparation) CreateFight(ctx context.Context, tgId string) (string, error) {

	code := s.codeGenerator.GetCode()

	if err := s.fight.Create(ctx, entity.Fight{
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

func (s battlePreparation) JoinFight(ctx context.Context, sessionID, tgId string) (entity.Fight, error) {

	preFight, err := s.fight.GetBySessionID(ctx, sessionID)
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
		s.helpers.setUserParams(&preFight.Users[i])

	}

	if err := s.fight.Update(ctx, preFight); err != nil {
		return entity.Fight{}, err
	}
	return preFight, nil
}
