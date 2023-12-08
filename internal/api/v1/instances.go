package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"github.com/zackarysantana/velocity/internal/workflows"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (a *V1App) PostInstanceStart() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.ParseBody(v1types.NewPostInstanceStartRequest),
		func(c *gin.Context) {
			data := middleware.GetBody(c).(v1types.PostInstanceStartRequest)

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

			projectId, err := primitive.ObjectIDFromHex(data.ProjectId)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error processing project id %v", err))
				return
			}

			workflow, err := data.Config.GetWorkflow(data.Workflow)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting workflow in config %v", err))
				return
			}

			jobs, err := workflows.GetJobsForWorkflow(*data.Config, workflow)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting jobs for instance %v", err))
				return
			}

			project, err := a.client.GetProject(c, projectId)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting project: %v", err))
				return
			}
			if project == nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("project not found %s", data.ProjectId))
				return
			}

			i := db.Instance{
				ProjectId: projectId,
				Config:    *data.Config,
			}
			instance, err := a.client.InsertInstance(c, &i)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error inserting instance %v", err))
				return
			}

			dbJobs := []*db.Job{}
			for _, job := range jobs {
				j := *job
				dbJobs = append(dbJobs, &db.Job{
					Name:          j.GetName(),
					Image:         j.GetImage(),
					Command:       j.GetCommand(),
					SetupCommands: j.GetSetupCommands(),
					Status:        jobtypes.JobStatusQueued,
					InstanceId:    &instance.Id,
				})
			}

			dbJobs, err = a.client.InsertJobs(c, dbJobs)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error inserting jobs %v", err))
				return
			}

			resp := v1types.PostInstanceStartResponse{
				InstanceId: instance.Id.Hex(),
				Jobs:       dbJobs,
			}
			c.JSON(200, resp)
		},
	}
}
