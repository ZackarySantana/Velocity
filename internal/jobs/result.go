package jobs

type JobResultSuccess struct {
	Logs string
}

type JobResultFailed struct {
	Error error
}

type JobResult struct {
	Job     Job
	Success *JobResultSuccess
	Failed  *JobResultFailed
}
