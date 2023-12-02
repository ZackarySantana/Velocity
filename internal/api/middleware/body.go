package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseBody(data interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		err := c.ShouldBindJSON(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Next()
	}
}
