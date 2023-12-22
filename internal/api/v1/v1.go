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

	a := V1App{client}

	// /api/v1
	v1 := router.Group("/api/v1")
	v1.POST("/first_time_register", a.PostFirstTimeRegister()...)

	// /api/v1/admin
	adminV1 := v1.Group("/admin", middleware.AdminAuth(client))
	adminV1.POST("/user", a.PostUser()...)

	authorizedV1 := v1.Group("/", middleware.Auth(client))

	// /api/v1/instances
	instances := authorizedV1.Group("/instances")
	instances.GET("/:instance_id", a.GetInstance()...)
	instances.POST("/start", a.PostInstanceStart()...)

	// /api/v1/jobs
	// TODO: Agent routes? should we separate these out?
	jobs := authorizedV1.Group("/jobs")
	jobs.POST("/dequeue", a.PostJobsDequeue()...)
	jobs.POST("/result", a.PostJobResult()...)

	return router, nil
}
