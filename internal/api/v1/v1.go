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

	v1 := router.Group("/api/v1")
	v1.POST("/first_time_register", a.PostFirstTimeRegister()...)

	adminV1 := v1.Group("/admin", middleware.AdminAuth(client))
	adminV1.POST("/user", a.PostUser()...)

	authorizedV1 := v1.Group("/", middleware.Auth(client))

	workflows := authorizedV1.Group("/workflows")
	workflows.POST("/start", a.PostInstanceStart()...)

	// TODO: Agent routes? should we separate these out?
	jobs := authorizedV1.Group("/jobs")
	jobs.POST("/dequeue", append(middleware.JobsFilter(postJobsDequeueOptsDefault), a.PostJobsDequeue)...)
	jobs.POST("/result", a.PostJobResult()...)

	return router, nil
}
