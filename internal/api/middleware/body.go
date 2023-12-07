package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ParseBody(newDataInstance func() interface{}) func(c *gin.Context) {
	return func(c *gin.Context) {
		data := newDataInstance()
		err := c.ShouldBindJSON(data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.Set("body", data)
		c.Next()
	}
}

func GetBody(c *gin.Context) interface{} {
	body, _ := c.Get("body")
	return body
}
