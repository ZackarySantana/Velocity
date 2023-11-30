package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/db"
)

func UseDB(client db.Connection) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Set("db", client)
		c.Next()
	}
}

func GetDB(c *gin.Context) db.Connection {
	d, exists := c.Get("db")
	if !exists {
		return db.Connection{}
	}
	return d.(db.Connection)
}
