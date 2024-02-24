package seabattle

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"os/signal"
	redisrep "seabattle/internal/repository/redis"
	"seabattle/internal/service"
	"syscall"

	"github.com/gin-gonic/gin"

	"seabattle/config"
	"seabattle/internal/transport/http/v1"
	"seabattle/pkg/client"
	"seabattle/pkg/httpserver"
	"seabattle/pkg/logger"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	ctx := context.Background()
	// Repository
	pg, err := client.NewPostgresClient(ctx, 5, cfg.PG)
	if err != nil {
		l.Fatal(fmt.Errorf("seabattle - Run - postgres.New: %w", err))
	}

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Address,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})
	defer pg.Close()
	repo := redisrep.New(clientRedis)

	// Service
	srvc := service.New(repo)

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, srvc)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("seabattle - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("seabattle - Run - httpServer.Notify: %w", err))

	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("seabattle - Run - httpServer.Shutdown: %w", err))
	}

}
