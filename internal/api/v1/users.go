package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func CreateUser() []gin.HandlerFunc {
	var data v1types.CreateUserRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"amount":  data,
			})
		},
	}
}
