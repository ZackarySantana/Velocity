package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/cli/logger"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/event"
	"github.com/zackarysantana/velocity/internal/event/meta"
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

func (a *Api) SendIndexesAppliedEvent(ctx *gin.Context) {
	user := middleware.MustGetAuthArtifact[db.User](ctx)
	err := a.es.SendEvent(ctx, event.Event{
		EventType: event.EventTypeIndexesApplied,
		Metadata:  meta.CreateApplyIndexes(user),
	})
	if err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePrivate,
		})
	}
}
