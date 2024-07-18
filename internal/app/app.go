package seabattle

import (
	"context"
	tgbotapi "github.com/go-telegram/bot"
	kafkago "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"os/signal"
	"seabattle/config"
	mongorep "seabattle/internal/repository/mongodb"
	"seabattle/internal/service"
	"seabattle/internal/transport/tg"
	"syscall"
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
	mong, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	collection := mong.Database("sea_battle").Collection("fight")

	rep := mongorep.New(collection)

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
	srvc := service.New(rep, gameCfg, kafka)

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
