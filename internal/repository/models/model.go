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

type User struct {
	SessionId        string `json:"tg_battle_session" redis:"tg_battle_session"`
	MyBattleField    string `json:"my_field" redis:"my_field"`
	EnemyBattleField string `json:"enemy_field" redis:"enemy_field"`
}
