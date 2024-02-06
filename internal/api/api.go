package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"github.com/zackarysantana/velocity/internal/db"
)

// Api is the main API server wrapper.
type Api struct {
	*gin.Engine

	client db.Database
}

func CreateApi(logger logger.Logger, client db.Database) *Api {
	api := Api{
		Engine: gin.New(),
		client: client,
	}
	api.Use(
		middleware.Logger(logger),
		gin.Recovery(),
		middleware.ErrorHandler(logger),
	)
	return &api
}

func (a *Api) AddUserRoutes() {
	user := a.Group("/user")
	user.Use(middleware.AuthUsernameAndPasswordUserWithMongoDB(a.client))
	user.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

func (a *Api) AddAgentRoutes() {
	agent := a.Group("/agent")
	agent.Use(middleware.AuthAgentWithMongoDB(a.client))
	agent.GET("/ping", func(c *gin.Context) {
		fmt.Println("TESTING")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
