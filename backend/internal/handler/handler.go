package handler

import (
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/repository"
)

type Handler struct {
    repo    *repository.RepoClient
    logger  *zap.Logger
}

func New(repo *repository.RepoClient, logger *zap.Logger) *Handler {
    return &Handler{
        repo:       repo,
        logger:     logger,
    }
}
