package seabattle

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram/bot"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"os/signal"
	"seabattle/config"
	"seabattle/internal/pb"
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

	fmt.Println(cfg.Mongo.URL)
	mong, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.URL))
	if err != nil {
		panic(err)
	}

	collection := mong.Database(cfg.Mongo.Db).Collection(cfg.Collection)

	fightRep := mongorep.NewFight(collection)
	userRep := mongorep.NewUser(collection)
	battlefieldRep := mongorep.NewBattleField(collection)

	gameCfg, err := config.NewGame("config/game.yaml")
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	battleActionSrvc := service.NewBattleAction(fightRep, gameCfg)
	battlePreparactionSrvc := service.NewBattlePreparation(fightRep, userRep, battlefieldRep, gameCfg)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	opts := []tgbotapi.Option{
		//tgbotapi.WithHTTPClient(time.Millisecond*20,httpserver.Port(cfg.Http.Port))

		//tgbotapi.WithMiddlewares(handlers.BreakSkat),
	}

	cc, err := grpc.NewClient("0.0.0.0:50501",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	grpcClient := pb.NewMatchmakingClient(cc)
	bot, err := tgbotapi.New(cfg.ApiToken, opts...)
	if err != nil {
		panic(err)
	}
	tg.New(bot, battleActionSrvc, battlePreparactionSrvc, grpcClient)
	bot.Start(ctx)

}
