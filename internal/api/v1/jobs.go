package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
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
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func (v *V1App) PostJobsEnqueue(c *gin.Context) {
	opts := middleware.GetJobsFilter(c)

	c.JSON(200, gin.H{
		"message": "pong",
		"amount":  opts.Amount,
	})
}
