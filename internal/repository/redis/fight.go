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

	if err := r.SetUser(ctx, fight.User1); err != nil {
		return err
	}
	if err := r.SetUser(ctx, fight.User2); err != nil {
		return err
	}

	return nil
}

func (r Redis) GetFight(ctx context.Context, sessionId string) (models.Fight, error) {
	var session models.Session
	if err := r.client.HGetAll(ctx, sessionId).Scan(&session); err != nil {
		return models.Fight{}, err
	}
	user1, err := r.GetUser(ctx, session.TgId1)
	if err != nil {
		return models.Fight{}, err
	}
	user2, err := r.GetUser(ctx, session.TgId2)
	if err != nil {
		return models.Fight{}, err
	}
	//user1.TgId = session.Turn

	fight := models.Fight{
		User1: user1,
		User2: user2,
		Turn:  session.Turn,
		State: session.Ready,
	}
	return fight, nil

}
