package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	postJobsDequeueOptsDefault = &middleware.JobFilterOpts{
		Amount: 5,
		Status: jobtypes.JobStatusQueued,
	}
)

func (v *V1App) PostJobsDequeue(c *gin.Context) {
	opts := middleware.GetJobsFilter(c)

	dbJobs, err := v.client.DequeueNJobs(c, int64(opts.Amount))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if len(dbJobs) == 0 {
		c.JSON(200, gin.H{
			"jobs": []interface{}{},
		})
		return
	}

	jobs := []db.Job{}
	for _, job := range dbJobs {
		jobs = append(jobs, *job)
	}

	resp := v1types.PostJobsDequeueResponse{Jobs: jobs}
	c.JSON(200, resp)
}

func (a *V1App) PostJobResult() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.ParseBody(v1types.NewPostJobResultRequest),
		func(c *gin.Context) {
			data := middleware.GetBody(c).(*v1types.PostJobResultRequest)
			id, err := primitive.ObjectIDFromHex(data.Id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			j := &db.Job{
				Id:     id,
				Status: jobtypes.JobStatusCompleted,
			}
			if data.Logs != nil {
				j.Logs = *data.Logs
			}
			if data.Error != nil {
				j.Error = *data.Error
			}
			fmt.Println("Updating job", j)
			j, err = a.client.UpdateJob(c, j)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			fmt.Println("Updated job", j)

			resp := v1types.PostJobResultResponse(*j)
			c.JSON(200, resp)
		},
	}
}
