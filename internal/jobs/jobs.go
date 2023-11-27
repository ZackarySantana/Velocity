package jobs

import "errors"

type Job struct {
	RepositoryURL string
	GitHash       string

	Command string
	Image   string

	Executor *JobExecutor
}

func (j *Job) Run() (string, error) {
	if j.Executor == nil {
		return "", errors.New("no executor provided")
	}
	return (*j.Executor).Execute(*j)
}
