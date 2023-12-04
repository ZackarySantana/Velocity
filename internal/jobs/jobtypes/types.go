package jobtypes

import "fmt"

type JobStatus string

var (
	JobStatusQueued    JobStatus = "queued"
	JobStatusInactive  JobStatus = "inactive"
	JobStatusActive    JobStatus = "active"
	JobStatusCompleted JobStatus = "completed"

	JobStatuses = []JobStatus{
		JobStatusQueued,
		JobStatusInactive,
		JobStatusActive,
		JobStatusCompleted,
	}
)

func JobStatusFromString(s string) (JobStatus, error) {
	for _, status := range JobStatuses {
		if s == string(status) {
			return status, nil
		}
	}
	return "", fmt.Errorf("invalid status %s", s)
}
