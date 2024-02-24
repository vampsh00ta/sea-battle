package main

import (
	"fmt"
	"seabattle/internal/service/entity"
)

func main() {
	// Configuration
	//cfg, err := config.New()
	//if err != nil {
	//	log.Fatalf("Config error: %s", err)
	//}
	ent := entity.NewBattleField()
	//ent.Fields[4][4].Ship = true
	//ent.Fields[4][5].Ship = true
	//
	//ent.Fields[4][4].Count = 2
	//ent.Fields[4][5].Count = 2
	var err error
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print(" ", ent.Fields[i][j].Marked)
		}
		fmt.Println()

	}
	ent.AddShip(entity.Point{4, 5}, entity.Point{4, 4}, 1)
	ent.Shoot(4, 5)
	ent.Shoot(4, 4)

	//
	//err = ent.Shoot(4, 4)
	//err = ent.Shoot(4, 5)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print(" ", ent.Fields[i][j].Marked)
		}
		fmt.Println()

	}
	fmt.Println(ent.Fields[0][0].Count)
	fmt.Println(err)

	// Run
	//seabattle.Run(cfg)
}
