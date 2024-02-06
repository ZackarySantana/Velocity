package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/cli/logger"
)

// Logger is a middleware that logs all requests.
func Logger(l logger.Logger) gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Output: l,
		Formatter: func(param gin.LogFormatterParams) string {
			if param.ErrorMessage == "" {
				param.ErrorMessage = "N/A"
			}
			param.ErrorMessage = strings.Trim(param.ErrorMessage, "\n")
			return fmt.Sprintf("[REQUEST] %s [%s] %s | %d | %s | %s | %s\n",
				param.TimeStamp.Format("2006-01-02 15:04:05"),
				param.Method,
				param.Path,
				param.StatusCode,
				param.ClientIP,
				param.Latency,
				param.ErrorMessage,
			)
		},
	})
}
