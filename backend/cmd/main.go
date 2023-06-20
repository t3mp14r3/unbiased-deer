package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/auth"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/config"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/gateway"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/handler"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/logger"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/nats"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/repository"
)

func main() {
    config := config.New()
    logger := logger.New()
    
    ctx, cancel := context.WithCancel(context.Background())

    repo := repository.New(ctx, &config.PostgresConfig, &config.RedisConfig, logger)
    defer repo.Close()
    
    auth := auth.New(ctx, &config.RedisConfig, logger)
    defer auth.Close()
    
    handler := handler.New(repo, logger)
    
    nats := nats.New(&config.NatsConfig, logger, handler)
    gateway := gateway.New(&config.GatewayConfig, nats, auth, logger)

    wg := &sync.WaitGroup{}
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := gateway.Run(ctx); err != nil {
            logger.Error("gateway run error", zap.Error(err))
            cancel()
        }
    }()

    logger.Info("gateway started", zap.String("addr", config.GatewayConfig.Addr))
    
    wg.Add(1)
    go func() {
        defer wg.Done()
        if err := nats.Subscribe(ctx); err != nil {
            logger.Error("nats consumer run error", zap.Error(err))
            cancel()
        }
    }()

    logger.Info("nats started")
    
    exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

    select {
        case <-ctx.Done():
            logger.Error("backend stop via context")
        case <-exit:
            logger.Info("backend stop")
    }

    cancel()
    wg.Wait()

    logger.Info("backend stopped")
}
