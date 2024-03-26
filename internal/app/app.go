package seabattle

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
	"log"
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

	gameCfg, err := config.NewGame()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	srvc := service.New(rep, psql, gameCfg)
	tr := tg.New(srvc)

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
