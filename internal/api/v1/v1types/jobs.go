package v1types

import "github.com/zackarysantana/velocity/internal/db"

type PostJobResultRequest struct {
	Id    string  `json:"id"`
	Logs  *string `json:"logs"`
	Error *string `json:"error"`
}

// TBA
type PostJobResultResponse db.Job

type PostJobsDequeueRequest struct{}

type PostJobsDequeueResponse struct {
	Jobs []db.Job `json:"jobs"`
}
