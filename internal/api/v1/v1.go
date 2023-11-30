package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/db"
)

func CreateV1App(client db.Connection) (*gin.Engine, error) {
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use((middleware.UseDB(client)))

	authorizedV1 := v1.Group("/")
	authorizedV1.Use(middleware.Auth)
	authorizedV1.GET("/jobs", middleware.QueryAmount(1), GetJobs)
	authorizedV1.GET("/jobs/dequeue", middleware.QueryAmount(1), GetJobs)

	return router, nil
}
