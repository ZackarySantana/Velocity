package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	apiv1 "github.com/zackarysantana/velocity/internal/api/v1"
	"github.com/zackarysantana/velocity/internal/db"
)

func main() {

	client, err := db.Connect(nil)
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	v1 := router.Group("/v1")
	v1.Use((middleware.UseDB(*client)))

	authorizedV1 := v1.Group("/")
	authorizedV1.Use(middleware.Auth)
	authorizedV1.GET("/jobs", middleware.QueryAmount(1), apiv1.GetJobs)
	authorizedV1.GET("/jobs/dequeue", middleware.QueryAmount(1), apiv1.GetJobs)

	// start
	router.Run(":8080")
}
