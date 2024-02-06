package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/cli/logger"
)

// ErrorHandler catches errors caused by middleware or handlers and logs them.
// The errors caught are from ctx.Error() calls. If the error is public, it will
// be returned to the client. If the error is private, it will be logged.
func ErrorHandler(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if len(c.Errors) > 0 {
			errMsg := ""
			privateErrLog := []string{}
			code := http.StatusInternalServerError
			for _, err := range c.Errors {
				if err.Type == gin.ErrorTypePublic {
					errMsg += err.Error()
				}
				if err.Type == gin.ErrorTypePrivate {
					privateErrLog = append(privateErrLog, err.Error())
				}
				if newCode, ok := err.Meta.(int); ok {
					code = newCode
				}
			}

			if errMsg == "" {
				errMsg = "an error occurred"
			}

			c.AbortWithStatusJSON(code, gin.H{"error": errMsg})

			if len(privateErrLog) > 0 {
				end := time.Now()
				latency := end.Sub(start)
				path := c.Request.URL.Path
				query := c.Request.URL.RawQuery
				clientIP := c.ClientIP()
				method := c.Request.Method

				logger.ErrorStr(fmt.Sprintf("%s [%s] %s?%s | %d | %s | %s | %s | %s\n", end.Format("2006-01-02 15:04:05"), method, path, query, code, clientIP, latency, strings.Join(privateErrLog, " | "), errMsg))
			}
		}
	}
}
