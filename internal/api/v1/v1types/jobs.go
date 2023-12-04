package v1types

import (
	"github.com/zackarysantana/velocity/internal/api/middleware"
	"github.com/zackarysantana/velocity/internal/db"
)

// POST /api/v1/jobs/result
type PostJobResultRequest struct {
	Id    string  `json:"id"`
	Logs  *string `json:"logs"`
	Error *string `json:"error"`
}
type PostJobResultResponse db.Job

// POST /api/v1/jobs/dequeue
type PostJobsDequeueRequest struct{}
type PostJobsDequeueQueryParams middleware.JobFilterOpts
type PostJobsDequeueResponse struct {
	Jobs []db.Job `json:"jobs"`
}
