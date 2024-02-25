package entity

import (
	"encoding/json"
	"errors"
	"seabattle/internal/repository/models"
	"seabattle/internal/service"
)

const (
	ship1 = iota
	ship2
	ship3
	ship4
)
const (
	heigth = 8
	weigth = 8
)

type BattleField struct {
	*models.BattleField
}

type Point struct {
	X, Y int
}

func NewBattleField() BattleField {
	fields := make([][]models.Field, heigth)
	for i := range fields {
		fields[i] = make([]models.Field, weigth)
	}
	return BattleField{
		&models.BattleField{
			fields,
		},
	}
}

func (b *BattleField) Shoot(x, y int) error {
	if b.Fields[x][y].Marked {
		return errors.New(service.AlreadyMarked)
	}
	if b.Fields[x][y].Ship {
		b.Fields[x][y].Marked = true

		used := make(map[Point]bool)

		var descCount func(x, y int)

		descCount = func(x, y int) {
			b.Fields[x][y].Count -= 1
			if b.Fields[x][y].Count == 0 {
				b.Fields[x][y].Dead = true
			}
			used[Point{x, y}] = true
			dirs := [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
			for _, dir := range dirs {
				x0, y0 := dir[0], dir[1]
				p := Point{x + x0, y + y0}
				if _, ok := used[p]; !ok && p.X < len(b.Fields[0]) && p.Y < len(b.Fields) && p.X >= 0 && p.Y >= 0 {
					if b.Fields[p.X][p.Y].Ship {
						descCount(p.X, p.Y)
					} else if b.Fields[x][y].Count == 0 {
						b.Fields[p.X][p.Y].Marked = true

					}
				}

			}

		}
		descCount(x, y)
		return nil
	} else {
		b.Fields[x][y].Marked = true
		return nil
	}

}
func (b *BattleField) AddShip(p1, p2 Point, shipType int) {
	x1 := min(p1.X, p2.X)
	x2 := max(p1.X, p2.X)
	y1 := min(p1.Y, p2.Y)
	y2 := max(p1.Y, p2.Y)
	if x1 == x2 {
		for i := y1; i <= y2; i++ {
			b.Fields[x1][i].Ship = true
			b.Fields[x1][i].Count = shipType + 1
		}
	} else {
		for i := x1; i <= x2; i++ {
			b.Fields[y1][i].Ship = true
			b.Fields[y1][i].Count = shipType + 1

		}
	}
}
func (b *BattleField) String() (string, error) {
	data, err := json.Marshal(b.BattleField)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

//func (b *BattleField) convertToModel() models.BattleField{
//	return models.BattleField{
//		b.Fields,
//	}
//}
