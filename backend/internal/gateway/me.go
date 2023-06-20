package gateway

import (
	"github.com/gin-gonic/gin"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (g *Gateway) me(c *gin.Context) {
    userID := c.GetString("userID")

    currency := c.Query("curr")

    if len(currency) == 0 {
        currency = "USD"
    }
    
    value := domain.Message{
        Type:       domain.MessageMe,
        Payload:    map[string]interface{}{
            "id": userID,
            "currency": currency,
        },
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
