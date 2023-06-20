package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/config"
)

type RepoClient struct {
    ctx     context.Context
    pc      *sqlx.DB
    rc      *redis.Client
    logger  *zap.Logger
}

func New(ctx context.Context, postgresConfig *config.PostgresConfig, redisConfig *config.RedisConfig, logger *zap.Logger) *RepoClient {
    conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		postgresConfig.Host,
		postgresConfig.Port,
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Name,
	)
    
    postgresConn, err := sqlx.Connect("postgres", conn)

	if err != nil {
        log.Fatalln("failed to initialize postgres connection! err:", err)
	}

    migrate(postgresConn.DB)

    redisConn := redis.NewClient(&redis.Options{
        Addr:     redisConfig.Addr,
        Password: redisConfig.Password,
        DB:       redisConfig.DB,
    })

    if err := redisConn.Ping(context.Background()).Err(); err != nil {
        log.Fatalln("failed to ping redis connection! err:", err)
    }

    return &RepoClient{
        ctx:    ctx,
        pc:     postgresConn,
        rc:     redisConn,
        logger: logger,
    }
}

func migrate(db *sql.DB) {
    if err := goose.SetDialect("postgres"); err != nil {
        log.Fatalln("failed to set goose dialect! err:", err)
    }

    if err := goose.Up(db, "migrations"); err != nil {
        log.Fatalln("failed to migrate the database! err:", err)
    }
}

func (r *RepoClient) Close() {
    if err := r.pc.Close(); err != nil {
        log.Fatalln("error while closing postgres connection! err:", err)
    }
   
    if err := r.rc.Close(); err != nil {
        log.Fatalln("error while closing redis connection! err:", err)
    }
}
