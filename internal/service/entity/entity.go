package entity

import (
	"errors"
	"fmt"
	"seabattle/internal/repository/models"
)

const (
	AlreadyMarkedErr  = "already_marked"
	AlreadyDeadErr    = "already_dead"
	AlreadySetErr     = "already_set"
	NotYourTurnErr    = "not_your_turn"
	GameEndedErr      = "game_ended"
	MaxShipCountErr   = "count_error"
	WrongPlacementErr = "wrong_placement"
)

const (
	maxShipCount = 4
	ShipType1    = iota
	ShipType2
	ShipType3
	ShipType4
)

const (
	Missed = iota
	Shooted
	Killed
	Lost = 3
)

const (
	StagePick = iota
	StageFight
	StageEnd
)

const (
	heigth = 8
	weigth = 8
)

type Point struct {
	X, Y int
}

func NewBattleField() *models.BattleField {
	fields := make([][]models.Field, heigth)
	for i := range fields {
		fields[i] = make([]models.Field, weigth)
	}
	return &models.BattleField{

		Fields: fields,
		Ships:  make(map[int]int),
	}
}

func Shoot(attacker, defender *models.BattleField, x, y int) (int, error) {
	res := Missed
	if attacker.Fields[x][y].Marked {
		res = -1
		return res, errors.New(AlreadyMarkedErr)
	}
	if defender.Fields[x][y].Ship {
		(*defender).Fields[x][y].Marked = true
		(*attacker).Fields[x][y].Marked = true
		(*attacker).Fields[x][y].Ship = true

		res = Shooted
		used := make(map[Point]bool)

		var descCount func(x, y int)

		descCount = func(x, y int) {
			(*defender).Fields[x][y].Count -= 1
			(*attacker).Fields[x][y].Count -= 1
			if defender.Fields[x][y].Count == 0 {
				(*defender).Fields[x][y].Dead = true
				(*attacker).Fields[x][y].Dead = true
				res = Killed
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
	if res == Killed {
		defender.Alive -= 1
	}
	if defender.Alive == 0 {
		res = Lost
	}
	fmt.Println(defender.Alive)
	return res, nil

}
func AddShip(b *models.BattleField, p1, p2 Point, shipType int) error {

	x1 := min(p1.X, p2.X)
	x2 := max(p1.X, p2.X)
	y1 := min(p1.Y, p2.Y)

	y2 := max(p1.Y, p2.Y)
	dirs := [][]int{{-1, 0}, {0, -1}, {1, 0}, {0, 1}}

	checkDir := func(x, y int, d [][]int) error {

		for _, dir := range d {
			x0, y0 := dir[0], dir[1]
			p := Point{x + x0, y + y0}

			if p.X < len(b.Fields[0]) && p.Y < len(b.Fields) && p.X >= 0 && p.Y >= 0 {
				if b.Fields[p.Y][p.X].Ship {
					return errors.New(WrongPlacementErr)
				}
			}
		}
		return nil
	}
	for i := y1; i >= y2; i-- {
	}
	if err := checkDir(x1, y1, dirs); err != nil {
		return err

	}
	if b.Ships[shipType] >= maxShipCount-shipType+1 {
		return errors.New(MaxShipCountErr)
	}
	if x1 == x2 {
		dirs = [][]int{{0, 1}, {1, 0}, {-1, 0}}
		for i := y1; i <= y2; i++ {
			if err := checkDir(x1, i, dirs); err != nil {
				for j := i; j >= y1; j-- {
					b.Fields[i][x1].Ship = false
					//b.Fields[x1][i].Marked = true
					b.Fields[i][x1].Count = 0
				}
				return err

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
				return err

			}
			b.Fields[y1][i].Ship = true
			//b.Fields[y1][i].Marked = true

			b.Fields[y1][i].Count = shipType

		}
	}
	b.Ships[shipType] += 1
	b.Alive += 1
	return nil

}

func IsReady(b *models.BattleField) bool {
	for i := range b.Ships {
		if b.Ships[i+1] < maxShipCount-i+1 {
			return false
		}
	}
	return true

}
