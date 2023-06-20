package gateway

import (
    "net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (g *Gateway) register(c *gin.Context) {
    var input domain.RegisterRequest

    if err := c.BindJSON(&input); err != nil {
        g.logger.Error("failed to parse request body", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": domain.ErrorBadBody,
        })
        return
    }
    
    if err := input.Validate(); err != nil {
        g.logger.Error("failed to validate request body", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    value := domain.Message{
        Type:       domain.MessageRegister,
        Payload:    map[string]interface{}{"name": input.Name},
    }

    resp := g.nats.Send(value)

    if resp.Err != nil {
        c.JSON(resp.Status(), gin.H{
            "error": resp.Err,
        })
        return
    }

    c.JSON(resp.Status(), resp.Data)
}
