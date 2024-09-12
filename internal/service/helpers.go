package service

import (
	"seabattle/config"
	"seabattle/internal/entity"
)

type helpers struct {
	gameConf *config.Game
}

type Helpers interface {
	setUserParams(user *entity.User)
	getCurrRoles(fight *entity.Fight) (attacker *entity.User, defender *entity.User)
}

func newHelpers(gameConf *config.Game) Helpers {
	return &helpers{
		gameConf: gameConf,
	}
}
func (s helpers) setUserParams(user *entity.User) {
	user.MyField = entity.NewBattleField(s.gameConf.Height, s.gameConf.Weight)
	user.EnemyField = entity.NewBattleField(s.gameConf.Height, s.gameConf.Weight)
	user.CurrX = -1
	user.CurrY = -1
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

func (s helpers) getCurrRoles(fight *entity.Fight) (attacker *entity.User, defender *entity.User) {
	if fight.Turn == fight.Users[0].TgId {
		return &fight.Users[0], &fight.Users[1]
	} else {
		return &fight.Users[1], &fight.Users[0]

	}

}
