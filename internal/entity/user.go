package entity

const (
	MyFieldQueryId    = "my_field_query_id"
	EnemyFieldQueryId = "enemy_field_query_id"
)

type User struct {
	MyField           *BattleField `bson:"my_field"`
	EnemyField        *BattleField `bson:"enemy_field"`
	TgId              string       `bson:"tg_id"`
	MyFieldQueryId    string       `json:"my_field_query_id" bson:"my_field_query_id"`
	EnemyFieldQueryId string       `json:"enemy_field_query_id" bson:"enemy_field_query_id"`
	CurrX             int          `json:"curr_x"  bson:"curr_x"`
	CurrY             int          `json:"curr_y"  bson:"curr_y"`
}
