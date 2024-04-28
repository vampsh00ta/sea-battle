package seabattle

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	"github.com/redis/go-redis/v9"
	kafkago "github.com/segmentio/kafka-go"
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

	kafka := &kafkago.Writer{
		Addr:                   kafkago.TCP("kafka:9092"),
		Topic:                  "search",
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
	srvc := service.New(rep, psql, gameCfg, kafka)

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

	tg.New(bot, srvc)

	bot.Start(ctx)

}
