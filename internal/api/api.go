package api

import (
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
