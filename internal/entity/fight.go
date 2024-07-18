package entity

type Fight struct {
	Turn      string `bson:"turn"`
	SessionId string `bson:"session_id"`
	Users     []User `bson:"users"`
	State     int    `bson:"state"`
	Ready     int    `bson:"ready"`
	Stage     int    `bson:"stage"`
}

type Shoot struct {
	Code  string `json:"code"`
	TgId  string `json:"tgId"`
	Point Point  `json:"point"`
}
type SetShip struct {
	Code  string `json:"code"`
	TgId  string `json:"tgId"`
	Point Point  `json:"point"`
}

type Point struct {
	X int `bson:"curr_x"`
	Y int `bson:"curr_y"`
}
type SearchFight struct {
	TgID   int `json:"tgID"`
	Rating int `json:"rating"`
}
