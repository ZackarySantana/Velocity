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

			// TODO: Add git repo / project to body/data
			// Find project in mongo by git repo
			// If not found, create project in mongo

			// Upload instance to mongo (connect it to the project)
			//  - Start workflow in mongo (logic should be here, calling client.InsertJobs)
			//  - Workflow is not it's own entity, in instance we should have a field for what workflow the instance is associated with, it's a 1-1 mapping

			// Example database calls:
			// ---
			// p, err := a.client.FindProject(projectFilter)
			// p, err := a.client.CreateProject(project)
			// ---
			// i, err := a.client.InsertInstance(instance)
			// ---
			// The jobs should all have a field that points back to the instance, like "instance_id"
			// We can make this optional or a only db field if we want to
			// j, err := a.client.InsertJobs(computeThoseJobs)
			// ---

			// Give user back instance id

			// resp := v1types.PostWorkflowsStartResponse(w.Id.Hex())
			// c.JSON(200, resp)
		},
	}
}
