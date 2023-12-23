package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *V1App) Health() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			err := a.client.Ping(c, nil)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			c.JSON(200, gin.H{
				"status": "ok",
			})
		},
	}
}
