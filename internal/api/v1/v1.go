package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/db"
)

type V1App struct {
	client db.Connection
}

func CreateV1App(client db.Connection) (*gin.Engine, error) {
	router := gin.Default()

	a := V1App{client: client}

	v1 := router.Group("/v1")

	authorizedV1 := v1.Group("/", middleware.Auth(client))
	authorizedV1.GET("/jobs", middleware.QueryAmount(1), a.GetJobs)
	authorizedV1.GET("/jobs/dequeue", middleware.QueryAmount(1), a.GetJobs)

	adminV1 := v1.Group("/admin", middleware.AdminAuth(client))
	adminV1.POST("/users", CreateUser()...)

	return router, nil
}
