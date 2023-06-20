package auth

import (
	"context"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/config"
)

type AuthClient struct {
    ctx     context.Context
    rc      *redis.Client
    logger  *zap.Logger
}

func New(ctx context.Context, redisConfig *config.RedisConfig, logger *zap.Logger) *AuthClient {
    redisConn := redis.NewClient(&redis.Options{
        Addr:     redisConfig.Addr,
        Password: redisConfig.Password,
        DB:       redisConfig.DB,
    })

    if err := redisConn.Ping(context.Background()).Err(); err != nil {
        log.Fatalln("failed to ping redis connection! err:", err)
    }

    return &AuthClient{
        ctx:    ctx,
        rc:     redisConn,
        logger: logger,
    }
}

func (a *AuthClient) Close() {
    if err := a.rc.Close(); err != nil {
        log.Fatalln("error while closing redis connection! err:", err)
    }
}

func (a *AuthClient) Auth(token string) (string, error) {
    userID, err := a.rc.Get(a.ctx, token).Result()
    
    if err != nil {
        a.logger.Error("failed to find user session record", zap.Error(err))
        return "", err
    }

    return userID, nil
}
