package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func (v *V1App) PostFirstTimeRegister() []gin.HandlerFunc {
	var data v1types.RegisterUserRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			anyAdminUsers, err := v.client.HasAdminUsers(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}
			if anyAdminUsers {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden- admin users already exist"})
				return
			}

			user, err := v.client.CreateAdminUser(c, data.Email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			c.JSON(200, user)
		},
	}
}
