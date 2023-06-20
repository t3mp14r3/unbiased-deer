package gateway

import (
	"context"
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/auth"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/config"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/nats"
)

type Gateway struct {
    r       *gin.Engine
    addr    string
    logger  *zap.Logger
    nats    *nats.NatsClient
    auth    *auth.AuthClient
}

func New(gatewayConfig *config.GatewayConfig, nats *nats.NatsClient, auth *auth.AuthClient, logger *zap.Logger) *Gateway {
    r := gin.Default()
    gateway := &Gateway{
        r:      r,
        addr:   gatewayConfig.Addr,
        logger: logger,
        nats:   nats,
        auth:   auth,
    }

    r.POST("/register", gateway.register)

    secure := r.Group("/")
    secure.Use(gateway.authMiddleware)

    secure.GET("/me", gateway.me)
    secure.POST("/deposit", gateway.deposit)
    secure.POST("/withdraw", gateway.withdraw)

    return gateway
}

func (g *Gateway) Run(ctx context.Context) error {
    errChan := make(chan error, 1)

    wg := &sync.WaitGroup{}
    wg.Add(1)

    go func() {
        defer wg.Done()
        if err := g.r.Run(g.addr); err != nil {
            g.logger.Error("gateway error", zap.Error(err))
            errChan <- err
        }
    }()

    var err error

    select {
        case <-ctx.Done():
            err = errors.New("gateway stop via context")
        case err = <-errChan:
    }

    wg.Wait()

    return err
}
