package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/cli/logger"
)

func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			errMsg := ""
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypePublic {
					errMsg += err.Error()
				}
				if err.Type == gin.ErrorTypePrivate {
					logger.Error(err)
				}
			}

			if errMsg == "" {
				errMsg = "an error occurred"
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": errMsg})
			c.Abort()
		}
	}
}
