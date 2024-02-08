package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
)

func (a *Api) AddUserRoutes() {
	user := a.Group("/user")
	user.Use(middleware.AuthUsernameAndPasswordUserWithMongoDB(a.db))
	user.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	superUser := user.Group("", middleware.OnlySuperUsers)
	superUser.POST("/create", a.CreateUser)
}
