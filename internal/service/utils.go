package service

import (
	"context"
	"seabattle/internal/entity"
	//"seabattle/internal/service/entity"
	"seabattle/internal/service/rules"
)

func (s service) setUserParams(user *entity.User) {
	user.MyField = entity.NewBattleField(s.gameConf.Height, s.gameConf.Weight)
	user.EnemyField = entity.NewBattleField(s.gameConf.Height, s.gameConf.Weight)
	user.CurrX = -1
	user.CurrY = -1
}

func (s service) setReady(ctx context.Context, res *int, ready int, token string) error {
	if ready != rules.PersonsReady {
		return nil
	}
	fight, err := s.mongo.GetFight(ctx, token)

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
	if err := s.mongo.UpdateFight(ctx, fight); err != nil {

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

func getCurrRoles(fight *entity.Fight) (attacker *entity.User, defender *entity.User) {
	if fight.Turn == fight.Users[0].TgId {
		return &fight.Users[0], &fight.Users[1]
	} else {
		return &fight.Users[1], &fight.Users[0]

	}

}
