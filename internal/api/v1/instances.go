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
			data := middleware.GetBody(c).(*v1types.PostInstanceStartRequest)

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

			var getProject func() (*db.Project, error)
			if data.ProjectId != nil {
				projectId, err := primitive.ObjectIDFromHex(*data.ProjectId)
				if err != nil {
					c.AbortWithStatusJSON(400, fmt.Sprint("error validating config: ", err))
					return
				}
				getProject = func() (*db.Project, error) {
					p, err := a.client.GetProjectById(c, projectId)
					if err != nil {
						return nil, fmt.Errorf("error getting project by id %v", err)
					}
					return p, nil
				}
			} else if data.ProjectName != nil {
				getProject = func() (*db.Project, error) {
					p, err := a.client.GetProjectByName(c, *data.ProjectName)
					if err != nil {
						return nil, fmt.Errorf("error getting project by name %s - %v", *data.ProjectName, err)
					}
					return p, nil
				}
			} else {
				c.AbortWithStatusJSON(400, "project name or id required")
				return
			}

			workflow, err := data.Config.GetWorkflow(data.Workflow)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting workflow in config %v", err))
				return
			}

			jobs, err := workflows.GetJobsForWorkflow(data.Config, workflow)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting jobs for instance %v", err))
				return
			}

			project, err := getProject()
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting project: %v", err))
				return
			}
			if project == nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("project not found %v", data.ProjectId))
				return
			}

			i := db.Instance{
				ProjectId: project.Id,
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
