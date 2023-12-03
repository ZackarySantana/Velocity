package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/api/middleware"
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

	c.JSON(200, gin.H{
		"message": "pong",
		"amount":  opts.Amount,
	})
}

func (v *V1App) PostJobsEnqueue(c *gin.Context) {
	opts := middleware.GetJobsFilter(c)

	c.JSON(200, gin.H{
		"message": "pong",
		"amount":  opts.Amount,
	})
}
