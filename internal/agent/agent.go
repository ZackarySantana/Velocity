package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
)

type Agent struct {
	Provider jobs.JobProvider
	Executor jobs.JobExecutor

	stop <-chan bool
	wg   *sync.WaitGroup
}

func NewAgent(provider jobs.JobProvider, executor jobs.JobExecutor, stop <-chan bool, wg *sync.WaitGroup) *Agent {
	return &Agent{provider, executor, stop, wg}
}

func (a *Agent) Start() error {
	limit := make(chan struct{}, 5)
	queue := make(chan jobs.Job)
	results := make(chan jobs.JobResult)

	go a.runJobs(queue, results, limit)
	go a.enqueue(queue, limit)
	go a.postResults(results)

	return nil
}

func (a *Agent) runJobs(queue <-chan jobs.Job, results chan<- jobs.JobResult, limit <-chan struct{}) {
	for job := range queue {
		a.wg.Add(1)
		go func(job jobs.Job) {
			defer func() {
				a.wg.Done()
				<-limit
			}()
			g := job.GetGitContext()
			ctx := jobs.NewContext(fmt.Sprintf("https://github.com/%s/%s", g.Owner, g.Repository), g.Hash)
			if g.URL != "" {
				ctx = jobs.NewContext(g.URL, g.Hash)
			}
			logs, err := a.Executor.Execute(ctx, job)
			job.SetStatus(jobtypes.JobStatusCompleted)
			fmt.Printf("Job complete (%d jobs running)\n", len(limit)-1)
			if err != nil {
				results <- jobs.JobResult{
					Job:    job,
					Failed: &jobs.JobResultFailed{Error: err},
				}
				return
			}
			results <- jobs.JobResult{
				Job:     job,
				Success: &jobs.JobResultSuccess{Logs: logs},
			}
		}(job)
	}
}

func (a *Agent) enqueue(queue chan<- jobs.Job, limit chan<- struct{}) {
	for {
		fmt.Println("Checking for jobs...")
		select {
		case <-a.stop:
			fmt.Println("Stopping agent...")
			return
		default:
		}
		jobs, err := a.Provider.Next(cap(limit) - len(limit))
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("Queuing %d jobs...\n", len(jobs))
		for _, job := range jobs {
			queue <- *job
			limit <- struct{}{}
		}

		finished, err := a.Provider.Finished()
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}

		if finished {
			fmt.Println("Finished queuing all jobs.")
			fmt.Println("Cleaning up...")
			err := a.Provider.Cleanup()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Finished cleanup.")
			fmt.Println(len(limit), "Job(s) are currently running.")
			return
		}

		time.Sleep(time.Second)
	}
}

func (a *Agent) postResults(results <-chan jobs.JobResult) {
	for result := range results {
		a.wg.Add(1)
		go func(result jobs.JobResult) {
			defer a.wg.Done()
			fmt.Println("Posting result...")
			err := a.Provider.Update(result)
			if err != nil {
				fmt.Println(err)
			}
		}(result)
	}
}
