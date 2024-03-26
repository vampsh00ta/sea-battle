package redis

import (
	"context"
	"seabattle/internal/repository/models"
	"sync"
)

type Fight interface {
	GetFight(ctx context.Context, sessionID string) (models.Fight, error)
	SetFight(ctx context.Context, fight models.Fight) error
}

func (r Redis) SetFight(ctx context.Context, fight models.Fight) error {
	wg := &sync.WaitGroup{}
	res := make(chan error, 2)
	users := []models.User{
		fight.User1,
		fight.User2,
	}
	for _, user := range users {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user models.User) {
			defer wg.Done()
			err := r.SetUser(ctx, user)
			res <- err

		}(wg, user)

	}

	wg.Wait()
	close(res)

	for err := range res {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Redis) GetFight(ctx context.Context, sessionId string) (models.Fight, error) {
	var session models.Session
	if err := r.client.HGetAll(ctx, sessionId).Scan(&session); err != nil {
		return models.Fight{}, err
	}

	wg := &sync.WaitGroup{}
	res := make(chan error, 2)
	tgIDs := []string{
		session.TgId1,
		session.TgId2,
	}
	resUsers := make([]models.User, 2)
	for i, tgID := range tgIDs {
		wg.Add(1)
		go func(wg *sync.WaitGroup, tgID string, i int) {
			defer wg.Done()
			user, err := r.GetUser(ctx, tgID)
			res <- err
			resUsers[i] = user

		}(wg, tgID, i)

	}

	wg.Wait()
	close(res)

	for err := range res {
		if err != nil {
			return models.Fight{}, err
		}
	}

	fight := models.Fight{
		User1: resUsers[0],
		User2: resUsers[1],
		Turn:  session.Turn,
		State: session.Ready,
	}
	return fight, nil

}
