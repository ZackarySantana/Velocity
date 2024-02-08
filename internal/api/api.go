package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/event"
)

// Api is the main API server wrapper.
type Api struct {
	*gin.Engine

	db db.Database
	es event.EventSender
}

func CreateApi(logger logger.Logger, db db.Database, es event.EventSender) *Api {
	api := Api{
		Engine: gin.New(),
		db:     db,
		es:     es,
	}
	api.Use(
		middleware.Logger(logger),
		gin.Recovery(),
		middleware.ErrorHandler(logger),
	)
	return &api
}
