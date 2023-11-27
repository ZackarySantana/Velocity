package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
)

func GetJobs(c *gin.Context) {
	amount := middleware.GetQueryAmount(c)

	c.JSON(200, gin.H{
		"message": "pong",
		"amount":  amount,
	})
}
