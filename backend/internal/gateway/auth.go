package gateway

import (
    "net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (g *Gateway) authMiddleware(c *gin.Context) {
    token := c.GetHeader("Authorization")

    if len(token) == 0 {
        g.logger.Error("unathorized request!")
        c.JSON(http.StatusUnauthorized, nil)
        c.Abort()
        return
    }

    userID, err := g.auth.Auth(c.GetHeader("Authorization"))

    if err != nil {
        g.logger.Error("unathorized request!", zap.Error(err))
        c.JSON(http.StatusUnauthorized, nil)
        c.Abort()
        return
    }

    c.Set("userID", userID)

    c.Next()
}
