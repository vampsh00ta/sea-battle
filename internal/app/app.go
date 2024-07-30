package seabattle

import (
	"context"
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
	//
	//clientRedis := redis.NewClient(&redis.Options{
	//	Addr:     cfg.Redis.Address,
	//	Password: cfg.Redis.Password,
	//	DB:       cfg.Redis.Db,
	//})
	//redrep := redisorep.New(clientRedis)
	//err := redrep.AddUserToQueue(ctx, 500, 555)

	mong, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Mongo.Address))
	if err != nil {
		panic(err)
	}

	collection := mong.Database(cfg.Mongo.Db).Collection(cfg.Collection)

	mongrep := mongorep.New(collection)

	gameCfg, err := config.NewGame()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Создаём соединение с gRPC-сервером SSO для клиента

	//kafka := &kafkago.Writer{
	//	Addr:                   kafkago.TCP("kafka:9092"),
	//	Topic:                  "search",
	//	Balancer:               &kafkago.LeastBytes{},
	//	AllowAutoTopicCreation: true,
	//}
	srvc := service.New(mongrep, gameCfg, nil)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	opts := []tgbotapi.Option{
		//tgbotapi.WithHTTPClient(time.Millisecond*20,httpserver.Port(cfg.Http.Port))

		//tgbotapi.WithMiddlewares(handlers.BreakSkat),
	}

	cc, err := grpc.DialContext(ctx, ":50501",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Создаём gRPC-клиент SSO/Auth
	grpcClient := pb.NewMatchmakingClient(cc)
	bot, err := tgbotapi.New(cfg.Apitoken, opts...)
	if err != nil {
		panic(err)
	}
	tg.New(bot, srvc, grpcClient)
	bot.Start(ctx)

}
