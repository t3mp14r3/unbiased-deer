package logger

import (
	"log"

	"go.uber.org/zap"
)

func New() *zap.Logger {
    logger, err := zap.NewProduction()

	if err != nil {
        log.Fatalln("failed to initialize zap logger! err:", err)
	}

    return logger
}
