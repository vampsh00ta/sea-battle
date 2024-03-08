package models

type User struct {
	TgId       string
	MyField    *BattleField
	EnemyField *BattleField
}
type Fight struct {
	Attacker  User
	Defender  User
	Turn      string
	SessionId string
	State     int
}
