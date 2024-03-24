package entity

import (
	"errors"
	"fmt"
	"math"
	"seabattle/internal/repository/models"
	"seabattle/internal/service/rules"
)

const (
	ShipType1 = iota
	ShipType2
	ShipType3
	ShipType4
)

type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func NewBattleField() *models.BattleField {
	fields := make([][]models.Field, rules.Height)
	for i := range fields {
		fields[i] = make([]models.Field, rules.Weight)
	}
	return &models.BattleField{

		Fields: fields,
		Ships:  make(map[int]int),
	}
}

func Shoot(attacker, defender *models.BattleField, x, y int) (int, error) {
	res := rules.Missed
	if attacker.Fields[x][y].Marked {
		res = -1
		return res, errors.New(rules.AlreadyMarkedErr)
	}
	if defender.Fields[x][y].Ship {
		(*defender).Fields[x][y].Marked = true
		(*defender).Fields[x][y].Shooted = true

		(*attacker).Fields[x][y].Marked = true
		(*attacker).Fields[x][y].Ship = true
		(*attacker).Fields[x][y].Shooted = true
		res = rules.Shooted
		used := make(map[Point]bool)

		var descCount func(x, y int)

		descCount = func(x, y int) {
			(*defender).Fields[x][y].Count -= 1
			(*attacker).Fields[x][y].Count -= 1
			if defender.Fields[x][y].Count == 0 {
				(*defender).Fields[x][y].Dead = true
				(*attacker).Fields[x][y].Dead = true
				res = rules.Killed
			}

			used[Point{x, y}] = true
			dirs := [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}
			for _, dir := range dirs {
				x0, y0 := dir[0], dir[1]
				p := Point{x + x0, y + y0}
				if _, ok := used[p]; !ok && p.X < len(defender.Fields[0]) && p.Y < len(defender.Fields) && p.X >= 0 && p.Y >= 0 {
					if defender.Fields[p.X][p.Y].Ship {
						descCount(p.X, p.Y)
					} else if defender.Fields[x][y].Count == 0 {
						(*defender).Fields[p.X][p.Y].Marked = true
						(*attacker).Fields[p.X][p.Y].Marked = true
					}
				}

			}

		}
		descCount(x, y)
	} else {
		defender.Fields[x][y].Marked = true
		attacker.Fields[x][y].Marked = true

	}
	if res == rules.Killed {
		defender.Alive -= 1
	}
	if defender.Alive == 0 {
		res = rules.Lost
	}
	fmt.Println(defender.Alive)
	return res, nil

}
func AddShip(b *models.BattleField, p1, p2 Point) (int, error) {

	x1 := min(p1.X, p2.X)
	x2 := max(p1.X, p2.X)
	y1 := min(p1.Y, p2.Y)

	y2 := max(p1.Y, p2.Y)
	dirs := [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 0}, {0, 1}}

	checkDir := func(x, y int, d [][]int) error {

		for _, dir := range d {
			x0, y0 := dir[0], dir[1]
			p := Point{x + x0, y + y0}

			if p.X < len(b.Fields[0]) && p.Y < len(b.Fields) && p.X >= 0 && p.Y >= 0 {
				if b.Fields[p.Y][p.X].Ship {
					return errors.New(rules.WrongPlacementErr)
				}
			}
		}
		return nil
	}

	if err := checkDir(x1, y1, dirs); err != nil {
		return 1, err

	}
	var shipType int
	shipType = int(math.Abs(float64((y2 - y1) + (x2 - x1))))
	if shipType >= 3 {
		return 0, errors.New(rules.WrongLengthErr)
	}
	if b.Ships[shipType] >= rules.MaxShipCount-shipType {
		return 0, errors.New(rules.MaxShipCountErr)
	}
	if x1 == x2 {

		dirs = [][]int{{0, 1}, {1, 0}, {-1, 0}}
		for i := y1; i <= y2; i++ {
			if err := checkDir(x1, i, dirs); err != nil {

				for j := i; j >= y1; j-- {
					b.Fields[j][x1].Ship = false
					//b.Fields[x1][i].Marked = true
					b.Fields[j][x1].Count = 0
				}

				return 0, err

			}
			b.Fields[i][x1].Ship = true
			//b.Fields[x1][i].Marked = true
			b.Fields[i][x1].Count = shipType + 1
		}
	} else {
		dirs = [][]int{{0, 1}, {0, -1}, {1, 0}}

		for i := x1; i <= x2; i++ {
			//fmt.Println(i)
			if err := checkDir(i, y1, dirs); err != nil {

				for j := i; j >= x1; j-- {
					b.Fields[y1][j].Ship = false
					//b.Fields[x1][i].Marked = true
					b.Fields[y1][j].Count = 0

				}
				return 0, err

			}
			b.Fields[y1][i].Ship = true
			//b.Fields[y1][i].Marked = true

			b.Fields[y1][i].Count = shipType + 1

		}
	}
	b.Ships[shipType] += 1
	b.Alive += 1

	var res int
	for t := range b.Ships {
		if b.Ships[t]+t == rules.MaxShipCount {
			res += 1
		}
	}
	if res == rules.ShipTypeCount {
		return rules.PersonsReady, nil
	}

	return 0, nil

}
