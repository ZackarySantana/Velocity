package main

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	apiv1 "github.com/zackarysantana/velocity/internal/api/v1"
)

func main() {
	// start server
	router := gin.Default()
	v1 := router.Group("/v1")

	v1.GET("/jobs", middleware.QueryAmount(5), apiv1.GetJobs)

	// start
	router.Run(":8080")
}
