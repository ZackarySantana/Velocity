package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func (v *V1App) PostUser() []gin.HandlerFunc {
	var data v1types.CreateUserRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			user, err := v.client.CreateUser(c, data.Email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			c.JSON(200, user)
		},
	}
}
