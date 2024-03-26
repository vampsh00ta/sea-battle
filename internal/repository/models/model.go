package models

type BattleField struct {
	Ships  map[int]int
	Alive  int       `json:"alive"`
	Fields [][]Field `json:"fields"`
}

type Field struct {
	Count int `json:"count"`
	Type  int `json:"type"`

	Ship    bool `json:"ship"`
	Marked  bool `json:"marked"`
	Dead    bool `json:"dead"`
	Shooted bool `json:"shooted"`
}

type Session struct {
	TgId1 string `json:"tgId1" redis:"tgId1"`
	TgId2 string `json:"tgId2" redis:"tgId2"`
	Turn  string `json:"turn" redis:"turn"`

	Ready int `json:"ready" redis:"ready"`
	Stage int `json:"stage" redis:"stage"`
}
