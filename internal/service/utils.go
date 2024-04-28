package service

import (
	"context"
	"seabattle/internal/models"
	//"seabattle/internal/service/entity"
	"seabattle/internal/service/rules"
)

func (s service) setUserParams(tgId string, user *models.User) {
	user.MyField = s.action.NewBattleField()
	user.EnemyField = s.action.NewBattleField()
	user.TgId = tgId
	user.CurrX = "-1"
	user.CurrY = "-1"
}

func (s service) setReady(ctx context.Context, res *int, ready int, token string) error {
	if ready != rules.PersonsReady {
		return nil
	}
	sessionId, err := s.psql.GetSessionByCode(ctx, token)
	if err != nil {
		return err
	}
	session, err := s.redis.GetSession(ctx, sessionId)
	if err != nil {
		return err
	}
	if ready == rules.PersonsReady {
		session.Ready += 1

	}
	switch session.Ready {
	case 1:
		*res = rules.PersonReady
	case 2:
		*res = rules.PersonsReady

	}
	if err := s.redis.SetSession(ctx, sessionId, session); err != nil {
		return err
	}

	return nil
}

//	func swapUsers(f) string {
//		var res string
//		fmt.Println(tgId1, tgId2, turn, "swap1")
//		if tgId1 == turn {
//			res = tgId2
//		} else {
//			res = tgId1
//		}
//		fmt.Println(tgId1, tgId2, turn, "swap2")
//
//		return res
//	}
func getAttacker(fight *models.Fight) {
	if fight.User1.TgId != fight.Turn {
		fight.User1, fight.User2 = fight.User2, fight.User1
	}
}
