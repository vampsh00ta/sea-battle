package models

type User struct {
	TgId              string
	MyField           *BattleField
	EnemyField        *BattleField
	CurrX             string `redis:"curr_x"`
	CurrY             string `redis:"curr_y"`
	MyFieldQueryId    string `redis:"my_field_query_id"`
	EnemyFieldQueryId string `redis:"enemy_field_query_id"`
}

type UserRedis struct {
	MyField           string `redis:"my_field"`
	EnemyField        string `redis:"enemy_field"`
	CurrX             string `redis:"curr_x"`
	CurrY             string `redis:"curr_y"`
	MyFieldQueryId    string `redis:"my_field_query_id"`
	EnemyFieldQueryId string `redis:"enemy_field_query_id"`
}

const (
	MyFieldQueryId    = "my_field_query_id"
	EnemyFieldQueryId = "enemy_field_query_id"
)

type Fight struct {
	User1     User
	User2     User
	Turn      string
	SessionId string
	State     int
}
