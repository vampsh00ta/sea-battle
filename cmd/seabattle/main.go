package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"seabattle/config"
	redisrep "seabattle/internal/repository/redis"
	"seabattle/internal/service/entity"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	ctx := context.Background()
	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	rep := redisrep.New(clientRedis)
	err = rep.CreateSessionByChatId(ctx, "key1", "key2")
	fmt.Println(err)

	fields := entity.NewBattleField()

	err = rep.SetBattleField(ctx, "key1", fields.BattleField, false)
	fmt.Println(err)

	err = rep.SetBattleField(ctx, "key1", fields.BattleField, true)
	fmt.Println(err)

	//res, err := rep.GetBattleField(ctx, "key1", false)
	//fmt.Println(res, err)

	user, err := rep.GetUserByChatId(ctx, "key1")
	fmt.Println(user, err)

	fields = entity.NewBattleField()

	err = rep.SetBattleField(ctx, "key2", fields.BattleField, false)
	fmt.Println(err)

	err = rep.SetBattleField(ctx, "key2", fields.BattleField, true)
	fmt.Println(err)

	//res, err := rep.GetBattleField(ctx, "key1", false)
	//fmt.Println(res, err)

	user, err = rep.GetUserByChatId(ctx, "key2")
	fmt.Println(user, err)

	//ent.Fields[4][4].Ship = true
	//ent.Fields[4][5].Ship = true
	//
	//ent.Fields[4][4].Count = 2
	//ent.Fields[4][5].Count = 2
	//var err error
	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		fmt.Print(" ", ent.Fields[i][j].Marked)
	//	}
	//	fmt.Println()
	//
	//}
	//ent.AddShip(entity.Point{4, 5}, entity.Point{4, 4}, 1)
	//ent.Shoot(4, 5)
	//ent.Shoot(4, 4)

	//
	//err = ent.Shoot(4, 4)
	//err = ent.Shoot(4, 5)

	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		fmt.Print(" ", ent.Fields[i][j].Marked)
	//	}
	//	fmt.Println()
	//
	//}
	//fmt.Println(ent.Fields[0][0].Count)
	//fmt.Println(err)

	// Run
	//seabattle.Run(cfg)
}
