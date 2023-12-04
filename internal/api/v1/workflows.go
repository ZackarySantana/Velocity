package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func (a *V1App) PostWorkflow() []gin.HandlerFunc {
	var data v1types.PostWorkflowRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			// post workflow to the database

			resp := v1types.PostWorkflowResponse("TBA")
			c.JSON(200, resp)
		},
	}
}
