package redis

import (
	"context"
	models2 "seabattle/internal/models"
	"sync"
)

type Fight interface {
	GetFight(ctx context.Context, sessionID string) (models2.Fight, error)
	SetFight(ctx context.Context, fight models2.Fight) error
}

func (r Redis) SetFight(ctx context.Context, fight models2.Fight) error {
	wg := &sync.WaitGroup{}
	res := make(chan error, 2)
	users := []models2.User{
		fight.User1,
		fight.User2,
	}
	for _, user := range users {
		wg.Add(1)
		go func(wg *sync.WaitGroup, user models2.User) {
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

func (r Redis) GetFight(ctx context.Context, sessionId string) (models2.Fight, error) {
	var session models2.Session
	if err := r.client.HGetAll(ctx, sessionId).Scan(&session); err != nil {
		return models2.Fight{}, err
	}

	wg := &sync.WaitGroup{}
	res := make(chan error, 2)
	tgIDs := []string{
		session.TgId1,
		session.TgId2,
	}
	resUsers := make([]models2.User, 2)
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
			return models2.Fight{}, err
		}
	}

	fight := models2.Fight{
		User1: resUsers[0],
		User2: resUsers[1],
		Turn:  session.Turn,
		State: session.Ready,
	}
	return fight, nil

}
