package models

type BattleField struct {
	Fields [][]Field `json:"fields"`
	Ships  map[int]int
	Alive  int `json:"alive"`
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
	Ready int    `json:"ready" redis:"ready"`
	Stage int    `json:"stage" redis:"stage"`
	Turn  string `json:"turn" redis:"turn"`
}
