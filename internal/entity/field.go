package entity

type BattleField struct {
	Ships  map[int]int `bson:"ships"`
	Fields [][]Field   `json:"fields" bson:"fields"`
	Alive  int         `json:"alive" bson:"alive"`
}

func NewBattleField(h, w int) *BattleField {
	fields := make([][]Field, h)
	for i := range fields {
		fields[i] = make([]Field, w)
	}
	return &BattleField{

		Fields: fields,
		Ships:  make(map[int]int),
	}
}

type Field struct {
	Count   int  `bson:"count"`
	Type    int  `bson:"type"`
	Ship    bool `bson:"ship"`
	Marked  bool `bson:"marked"`
	Dead    bool `bson:"dead"`
	Shooted bool `bson:"shooted"`
}
