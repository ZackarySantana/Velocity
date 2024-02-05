package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Api struct {
	*gin.Engine
	client *mongo.Client
}

func CreateApi(logger logger.Logger, client *mongo.Client) *Api {
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

func (a *Api) AddAgentRoutes(db, collection string) {
	agent := a.Group("/agent")
	agent.Use(middleware.AuthWithMongoDBAndUsernameAndPasswordFromJSONBody(a.client, db, collection))
	agent.GET("/ping", func(c *gin.Context) {
		fmt.Println("TESTING")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}
