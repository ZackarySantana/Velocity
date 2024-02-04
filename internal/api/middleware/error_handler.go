package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/cli/logger"
)

func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			errMsg := ""
			privateErrLog := []string{}
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypePublic {
					errMsg += err.Error()
				}
				if err.Type == gin.ErrorTypePrivate {
					privateErrLog = append(privateErrLog, err.Error())
				}
			}

			if errMsg == "" {
				errMsg = "an error occurred"
			}

			if len(privateErrLog) > 0 {
				logger.WrapError(strings.Join(privateErrLog, " | "))
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errMsg})
		}
	}
}
