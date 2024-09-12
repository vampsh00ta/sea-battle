package service

import (
	"context"
	"errors"
	"seabattle/config"
	irep "seabattle/internal/app/repository/mongorep"
	isrvc "seabattle/internal/app/service"
	"seabattle/internal/entity"
	"seabattle/internal/service/rules"
)

type battleAction struct {
	fight   irep.Fight
	entity  EntityInteraction
	helpers Helpers
}

func NewBattleAction(
	fightRepo irep.Fight,

	game *config.Game) isrvc.BattleAction {

	return &battleAction{
		fight:   fightRepo,
		entity:  newEntityInteraction(game),
		helpers: newHelpers(game),
	}
}

func (s battleAction) Shoot(ctx context.Context, req entity.Shoot) (entity.Fight, int, error) {

	fight, err := s.fight.GetBySessionID(ctx, req.Code)
	if err != nil {
		return entity.Fight{}, -1, err
	}
	switch {
	case fight.Stage == rules.StageEnd:
		return entity.Fight{}, -1, errors.New(rules.GameEndedErr)
	case fight.Turn != req.TgId:
		return entity.Fight{}, -1, errors.New(rules.NotYourTurnErr)
	}

	attacker, defender := s.helpers.getCurrRoles(&fight)
	res, err := s.entity.shoot(attacker.EnemyField, defender.MyField, req.Point.Y, req.Point.X)
	if err != nil {
		return entity.Fight{}, -1, err
	}
	switch res {
	case rules.Missed:
		fight.Turn = defender.TgId
	case rules.Lost:
		fight.State = rules.StageEnd

	}

	if err := s.fight.Update(ctx, fight); err != nil {
		return entity.Fight{}, -1, err
	}

	return fight, res, nil

}
