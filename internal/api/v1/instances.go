package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"github.com/zackarysantana/velocity/internal/workflows"
	"github.com/zackarysantana/velocity/src/clients/v1types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Get instance
func (a *V1App) GetInstance() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.InstanceId(),
		func(c *gin.Context) {
			instance_id := middleware.GetInstanceId(c)

			instance, err := a.client.GetInstanceById(c, instance_id)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting instance %v", err))
				return
			}

			jobs, err := a.client.GetJobsByInstanceId(c, instance_id)
			if err != nil {
				c.AbortWithStatusJSON(400, fmt.Sprintf("error getting jobs for instance %v", err))
				return
			}

			j := []db.Job{}
			for _, job := range jobs {
				j = append(j, *job)
			}
			res := v1types.GetInstanceResponse{
				Instance: *instance,
				Jobs:     j,
			}
			c.JSON(200, res)
		},
	}
}

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
