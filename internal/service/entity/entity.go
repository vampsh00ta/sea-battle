package entity

import (
	"errors"
	"seabattle/internal/service"
)

// * -уничтоженный
// . - занятый
const (
	ship1 = iota
	ship2
	ship3
	ship4
)

type BattleField struct {
	Fields [][]Field
}
type Field struct {
	Count  int
	Ship   bool
	Marked bool
	Dead   bool
}
type Point struct {
	X, Y int
}

func NewBattleField() BattleField {
	fields := make([][]Field, 8)
	for i := range fields {
		fields[i] = make([]Field, 8)
	}
	return BattleField{
		fields,
	}
}

func (b *BattleField) Shoot(x, y int) error {
	if b.Fields[x][y].Marked {
		return errors.New(service.AlreadyMarked)
	}
	if b.Fields[x][y].Ship {
		b.Fields[x][y].Marked = true
		//b.Fields[x][y].Count -= 1
		//if b.Fields[x][y].Count == 0 {
		//
		//	used := make(map[point]bool)
		//	var fillDeadShip func(x, y int)
		//	fillDeadShip = func(x, y int) {
		//		fmt.Println(x, y)
		//		dirs := [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
		//		for _, dir := range dirs {
		//			x0, y0 := dir[0], dir[1]
		//			p := point{x + x0, y + y0}
		//
		//			if _, ok := used[p]; !ok && p.x < len(b.Fields[0]) && p.y < len(b.Fields) && p.x >= 0 && p.y >= 0 {
		//				if b.Fields[p.x][p.y].Ship {
		//					used[p] = true
		//					fillDeadShip(p.x, p.y)
		//				} else {
		//					fmt.Println(p)
		//
		//					b.Fields[p.x][p.y].Marked = true
		//				}
		//			}
		//
		//		}
		//
		//	}
		//	fillDeadShip(x, y)
		//
		//} else {
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
		//}
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

//func (b *BattleField) convertToModel() models.BattleField{
//	return models.BattleField{
//		b.Fields,
//	}
//}
