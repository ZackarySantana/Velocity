package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
)

func (a *V1App) PostWorkflowsStart() []gin.HandlerFunc {
	var data v1types.PostWorkflowsStartRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			err := data.Config.Populate()
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error populating config: %v", err))
				return
			}
			err = data.Config.Validate()
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error validating config: %v", err))
				return
			}

			// Find project in mongo by git repo
			// If not found, create project in mongo

			// Upload instance to mongo
			// Start workflow in mongo (logic should be here, calling client.InsertJobs)
			// w, err := a.client.StartWorkflow(data.Config, data.Workflow)
			// j, err := a.client.InsertJobs(computeThoseJobs)

			// Give user back instance id

			// resp := v1types.PostWorkflowsStartResponse(w.Id.Hex())
			// c.JSON(200, resp)
		},
	}
}
