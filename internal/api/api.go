package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/cli/logger"
)

// TODO: Should logger be the best way to handle logs here?
// should every route get access to this logger?
func ErrorHandler(logger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error(err)
			}

			c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred"})
			c.Abort()
		}
	}
}
