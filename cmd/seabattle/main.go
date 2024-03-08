package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"seabattle/config"
	redisrep "seabattle/internal/repository/redis"
	"seabattle/internal/service"
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
	//id, err := rep.CreateSessionByChatId(ctx, "key1", "key2")
	if err != nil {
		panic(err)
	}
	//session := models.Session{
	//	TgId1: "key1",
	//	TgId2: "key2",
	//	Ready: 0,
	//	Stage: models.StageFight,
	//	Turn:  "key1",
	//}
	//if err := rep.SetSession(ctx, id, session); err != nil {
	//	panic(err)
	//}
	//bf1 := entity.NewBattleField()
	//bf2 := entity.NewBattleField()
	//bf3 := entity.NewBattleField()
	//bf4 := entity.NewBattleField()
	//entity.AddShip(bf3, entity.Point{X: 1, Y: 5}, entity.Point{X: 6, Y: 5}, 5)

	//user1 := models.User{"key1", bf1, bf2}
	//user2 := models.User{"key2", bf3, bf4}

	p := entity.Point{X: 3, Y: 5}

	srvc := service.New(rep)
	f, err := srvc.InitSessionFight(ctx, "key1", "key2")
	fmt.Println(f)
	if err != nil {
		panic(err)
	}
	fight := entity.Fight{"key1", "key2", "key1", f.SessionId, -1}

	err = srvc.AddShip(ctx, "key2", entity.Point{X: 1, Y: 5}, entity.Point{X: 5, Y: 5}, entity.ShipType4)
	fmt.Println(err)
	err = srvc.AddShip(ctx, "key2", entity.Point{X: 7, Y: 3}, entity.Point{X: 7, Y: 6}, entity.ShipType3)
	fmt.Println(err)
	b, _ := rep.GetBattleField(ctx, "key2", true)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print("|", b.Fields[i][j].Marked, b.Fields[i][j].Ship, " ")
		}
		fmt.Println()
	}
	a := entity.IsReady(b)
	fmt.Println(a)
	fmt.Println(err)
	_, res, err := srvc.Shoot(ctx, fight, entity.Point{1, p.Y})
	fmt.Println(res, err)
	_, res, err = srvc.Shoot(ctx, fight, entity.Point{2, p.Y})
	fmt.Println(res, err)

	_, res, err = srvc.Shoot(ctx, fight, entity.Point{3, p.Y})
	fmt.Println(res, err)

	_, res, err = srvc.Shoot(ctx, fight, entity.Point{4, p.Y})
	fmt.Println(res, err)

	_, res, err = srvc.Shoot(ctx, fight, entity.Point{5, p.Y})
	fmt.Println(res, err)
	_, res, err = srvc.Shoot(ctx, fight, entity.Point{6, p.Y})
	fmt.Println(res, err)

	_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 3})
	fmt.Println(res, err)
	_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 2})
	fmt.Println(res, err)
	_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 2})
	fmt.Println(res, err)
	b1, _ := rep.GetBattleField(ctx, "key2", true)
	b2, _ := rep.GetBattleField(ctx, "key1", false)

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print("|", b1.Fields[i][j].Marked, b1.Fields[i][j].Ship, " ")
		}
		fmt.Println()
	}
	fmt.Println()
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Print("|", b2.Fields[i][j].Marked, b2.Fields[i][j].Ship, " ")
		}
		fmt.Println()
	}
	//fields := entity.NewBattleField()
	//
	//err = rep.SetBattleField(ctx, "key1", fields.BattleField, false)
	//fmt.Println(err)
	//
	//err = rep.SetBattleField(ctx, "key1", fields.BattleField, true)
	//fmt.Println(err)
	//
	////res, err := rep.GetBattleField(ctx, "key1", false)
	////fmt.Println(res, err)
	//
	//user, err := rep.GetUserByChatId(ctx, "key1")
	//fmt.Println(user, err)
	//
	//fields = entity.NewBattleField()
	//
	//err = rep.SetBattleField(ctx, "key2", fields.BattleField, false)
	//fmt.Println(err)
	//
	//err = rep.SetBattleField(ctx, "key2", fields.BattleField, true)
	//fmt.Println(err)
	//
	////res, err := rep.GetBattleField(ctx, "key1", false)
	////fmt.Println(res, err)
	//
	//user, err = rep.GetUserByChatId(ctx, "key2")
	//fmt.Println(user, err)

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
