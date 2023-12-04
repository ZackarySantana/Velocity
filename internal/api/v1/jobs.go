package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/api/v1/v1types"
	"github.com/zackarysantana/velocity/internal/db"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	getJobsOptsDefault = &middleware.JobFilterOpts{
		Amount: 100,
		Status: jobtypes.JobStatusCompleted,
	}

	postJobsDequeueOptsDefault = &middleware.JobFilterOpts{
		Amount: 5,
		Status: jobtypes.JobStatusQueued,
	}
)

func (v *V1App) GetJobs(c *gin.Context) {
	opts := middleware.GetJobsFilter(c)

	c.JSON(200, gin.H{
		"message": "pong",
		"amount":  opts.Amount,
		"status":  opts.Status,
	})
}

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
	c.JSON(200, gin.H{
		"jobs": dbJobs,
	})
}

func (a *V1App) PostJobResult() []gin.HandlerFunc {
	var data v1types.PostJobResultRequest
	return []gin.HandlerFunc{
		middleware.ParseBody(&data),
		func(c *gin.Context) {
			id, err := primitive.ObjectIDFromHex(data.Id)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			j := &db.Job{
				Id:     id,
				Status: jobtypes.JobStatusCompleted,
			}
			j, err = a.client.UpdateJob(c, j)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
				return
			}

			resp := v1types.PostJobResultResponse(*j)
			c.JSON(200, resp)
		},
	}
}
