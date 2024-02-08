package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
)

func (a *Api) AddAgentRoutes() {
	agent := a.Group("/agent")
	agent.Use(middleware.AuthAgentWithMongoDB(a.db))
	agent.GET("/ping", func(c *gin.Context) {
		fmt.Println("TESTING")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
