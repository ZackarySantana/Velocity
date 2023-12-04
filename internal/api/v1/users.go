package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func (v *V1App) PostUser() []gin.HandlerFunc {
	var data v1types.PostUserRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			user, err := v.client.InsertUser(c, data.Email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			resp := v1types.PostUserResponse(*user)
			c.JSON(200, resp)
		},
	}
}
