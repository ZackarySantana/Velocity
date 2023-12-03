package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zackarysantana/velocity/internal/jobs"
)

type JobFilterOpts struct {
	Amount int
	Status jobs.JobStatus
}

var (
	defaultJobFilter = JobFilterOpts{
		Amount: 100,
		Status: jobs.JobStatusQueued,
	}
)

func JobsFilter(opts *JobFilterOpts) []gin.HandlerFunc {
	if opts == nil {
		opts = &defaultJobFilter
	}
	return []gin.HandlerFunc{
		func(c *gin.Context) {
			amountAsString := c.DefaultQuery("amount", strconv.Itoa(opts.Amount))
			amount, err := strconv.Atoi(amountAsString)
			if err != nil {
				c.JSON(400, gin.H{
					"message": "amount must be a number. received: " + amountAsString,
				})
				c.Abort()
				return
			}
			statusAsString := c.DefaultQuery("status", string(opts.Status))
			status, err := jobs.JobStatusFromString(statusAsString)
			if err != nil {
				c.JSON(400, gin.H{
					"message": "status must be a valid job status. received: " + statusAsString,
				})
				c.Abort()
				return
			}

			c.Set("amount", amount)
			c.Set("status", status)
			c.Next()
		},
	}
}

func GetJobsFilter(c *gin.Context) JobFilterOpts {
	amount, aExists := c.Get("amount")
	status, sExists := c.Get("status")
	if !aExists || !sExists {
		return JobFilterOpts{
			Amount: -1,
			Status: jobs.JobStatusCompleted,
		}
	}
	return JobFilterOpts{
		Amount: amount.(int),
		Status: status.(jobs.JobStatus),
	}
}
