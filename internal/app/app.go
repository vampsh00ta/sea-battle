package seabattle

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	psqlrep "seabattle/internal/repository/psql"
	redisrep "seabattle/internal/repository/redis"
	"seabattle/internal/service"
	"seabattle/internal/transport/tg"
	"syscall"

	"seabattle/config"
	"seabattle/pkg/client"
)

//// Run creates objects via constructors.
//func Run(cfg *config.Config) {
//	l := logger.New(cfg.Log.Level)
//
//	ctx := context.Background()
//	// Repository
//
//
//	// HTTP Server
//	handler := gin.New()
//	v1.NewRouter(handler, l, srvc)
//	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
//
//	// Waiting signal
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
//
//	select {
//	case s := <-interrupt:
//		l.Info("seabattle - Run - signal: " + s.String())
//	case err = <-httpServer.Notify():
//		l.Error(fmt.Errorf("seabattle - Run - httpServer.Notify: %w", err))
//
//	}
//
//	// Shutdown
//	err = httpServer.Shutdown()
//	if err != nil {
//		l.Error(fmt.Errorf("seabattle - Run - httpServer.Shutdown: %w", err))
//	}
//
//}

func NewPooling(cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	//ctx := context.Background()
	pg, err := client.NewPostgresClient(ctx, 5, cfg.PG)
	if err != nil {
		panic(err)
	}
	defer pg.Close()
	psql := psqlrep.New(pg)

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
	rep := redisrep.New(clientRedis)

	//// Service
	//p := entity.Point{X: 3, Y: 5}
	//
	srvc := service.New(rep, psql)
	tr := tg.New(srvc)
	//
	//f, err := srvc.InitSessionFight(ctx, "key1", "key2")
	//fmt.Println(f)
	//if err != nil {
	//	panic(err)
	//}
	//fight := entity.Fight{"key1", "key2", "key1", f.SessionId, -1}
	//
	//err = srvc.AddShip(ctx, "key2", entity.Point{X: 1, Y: 5}, entity.Point{X: 5, Y: 5}, entity.ShipType4)
	//fmt.Println(err)
	//err = srvc.AddShip(ctx, "key2", entity.Point{X: 7, Y: 3}, entity.Point{X: 7, Y: 6}, entity.ShipType3)
	//fmt.Println(err)
	//b, _ := rep.GetBattleField(ctx, "key2", true)
	//
	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		fmt.Print("|", b.Fields[i][j].Marked, b.Fields[i][j].Ship, " ")
	//	}
	//	fmt.Println()
	//}
	//a := entity.IsReady(b)
	//fmt.Println(a)
	//fmt.Println(err)
	//_, res, err := srvc.Shoot(ctx, fight, entity.Point{1, p.Y})
	//fmt.Println(res, err)
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{2, p.Y})
	//fmt.Println(res, err)
	//
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{3, p.Y})
	//fmt.Println(res, err)
	//
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{4, p.Y})
	//fmt.Println(res, err)
	//
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{5, p.Y})
	//fmt.Println(res, err)
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{6, p.Y})
	//fmt.Println(res, err)
	//
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 3})
	//fmt.Println(res, err)
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 2})
	//fmt.Println(res, err)
	//_, res, err = srvc.Shoot(ctx, fight, entity.Point{7, 2})
	//fmt.Println(res, err)
	//b1, _ := rep.GetBattleField(ctx, "key2", true)
	//b2, _ := rep.GetBattleField(ctx, "key1", false)
	//
	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		fmt.Print("|", b1.Fields[i][j].Marked, b1.Fields[i][j].Ship, " ")
	//	}
	//	fmt.Println()
	//}
	//fmt.Println()
	//for i := 0; i < 8; i++ {
	//	for j := 0; j < 8; j++ {
	//		fmt.Print("|", b2.Fields[i][j].Marked, b2.Fields[i][j].Ship, " ")
	//	}
	//	fmt.Println()
	//}

	//res, err := rep.GetBattleField(ctx, "key1", false)
	//fmt.Println(res, err)

	//app.NewPooling(cfg)

	//Run
	//seabattle.Run(cfg)
	//db, err := client.NewPostgresClient(ctx, 5, cfg.PG)
	//if err != nil {
	//	panic(err)
	//}
	//if err != nil {
	//	panic(err)
	//}
	//rep := rep.New(db)
	//fmt.Println(res, a)
	//tx := rep.GetDb()
	//a, err := rep.GetAllSubjectsOrderByName(ctx, true)
	//fmt.Println(a, err)
	//auth := &authentication.AuthMap{DB: make(map[int64]*authentication.User)}
	//auth.LogIn(564764193, 955, 2)

	//if err != nil {
	//	panic(err)
	//}
	//srvc := service.New(rep)
	//log := logger.New(cfg.Level)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	opts := []tgbotapi.Option{
		//tgbotapi.WithHTTPClient(time.Millisecond*20,httpserver.Port(cfg.Http.Port))

		//tgbotapi.WithMiddlewares(handlers.BreakSkat),
	}

	bot, err := tgbotapi.New(cfg.Apitoken, opts...)
	if err != nil {
		panic(err)
	}
	//handlers.New(bot, srvc, log)
	//bot.DeleteWebhook(ctx, &tgbotapi.DeleteWebhookParams{
	//	true,
	//})

	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/creategame",
		tgbotapi.MatchTypeExact, tr.CreateFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/start",
		tgbotapi.MatchTypePrefix, tr.JoinFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeMessageText,
		"/joingame",
		tgbotapi.MatchTypePrefix, tr.JoinFight,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"set#",
		tgbotapi.MatchTypePrefix, tr.SetShip,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"pass#",
		tgbotapi.MatchTypePrefix, tr.Pass,
	)
	bot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData,
		"shoot#",
		tgbotapi.MatchTypePrefix, tr.GameAction,
	)

	bot.Start(ctx)

}
