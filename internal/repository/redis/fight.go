package redis

import (
	"context"
	"seabattle/internal/repository/models"
)

type Fight interface {
	GetFight(ctx context.Context, sessionID string) (models.Fight, error)
	SetFight(ctx context.Context, fight models.Fight) error
}

func (r Redis) SetFight(ctx context.Context, fight models.Fight) error {

	if err := r.SetBattleField(ctx, fight.Attacker.TgId, fight.Attacker.MyField, true); err != nil {
		return err
	}
	if err := r.SetBattleField(ctx, fight.Attacker.TgId, fight.Attacker.EnemyField, false); err != nil {
		return err
	}
	if err := r.SetBattleField(ctx, fight.Defender.TgId, fight.Defender.MyField, true); err != nil {
		return err
	}
	if err := r.SetBattleField(ctx, fight.Defender.TgId, fight.Defender.EnemyField, false); err != nil {
		return err
	}
	return nil
}
func (r Redis) GetFight(ctx context.Context, sessionId string) (models.Fight, error) {
	var session models.Session
	var err error
	if err := r.client.HGetAll(ctx, sessionId).Scan(&session); err != nil {
		return models.Fight{}, err
	}
	var attacker, defender models.User
	attacker.TgId = session.Turn
	f1, err := r.GetBattleFields(ctx, attacker.TgId)
	if err != nil {
		return models.Fight{}, err
	}
	attacker.MyField = f1[0]
	attacker.EnemyField = f1[1]

	if session.TgId2 != session.Turn {
		defender.TgId = session.TgId2
	} else {
		defender.TgId = session.TgId1
	}
	f2, err := r.GetBattleFields(ctx, defender.TgId)
	if err != nil {
		return models.Fight{}, err
	}
	defender.MyField = f2[0]
	defender.EnemyField = f2[1]
	fight := models.Fight{
		Attacker: attacker,
		Defender: defender,
		Turn:     session.Turn,
		State:    session.Ready,
	}
	return fight, nil

}
