package models

type BattleField struct {
	Fields [][]Field `json:"fields"`
}

type Field struct {
	Count  int  `json:"count"`
	Ship   bool `json:"ship"`
	Marked bool `json:"marked"`
	Dead   bool `json:"dead"`
}

const (
	BattleSession = "tg_battle_session"
	MyField       = "my_field"
	EnemyField    = "enemy_field"
)

const (
	StagePicking = "pick"
	StageFight   = "fight"
)

type Session struct {
	TgId1 string `json:"tgId1" redis:"tgId1"`
	TgId2 string `json:"tgId2" redis:"tgId2"`
	Ready int    `json:"ready" redis:"ready"`
	Stage string `json:"stage" redis:"stage"`
	Step  string `json:"step" redis:"step"`
}

type User struct {
	MyBattleField    string `json:"my_field" redis:"my_field"`
	EnemyBattleField string `json:"enemy_field" redis:"enemy_field"`
}
