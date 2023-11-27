package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/jobs"
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
			job.Executor = &a.Executor
			logs, err := job.Run()
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
		default:
		}
		jobs, err := a.Provider.Next(cap(limit) - len(limit))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf("Queuing %d jobs...\n", len(jobs))
		for _, job := range jobs {
			queue <- job
			limit <- struct{}{}
		}

		time.Sleep(time.Second)
	}
}

func (a *Agent) postResults(results <-chan jobs.JobResult) {
	for result := range results {
		err := a.Provider.Update(result)
		if err != nil {
			fmt.Println(err)
		}
	}
}
