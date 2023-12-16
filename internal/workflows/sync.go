package workflows

import (
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/src/config"
)

func RunSyncWorkflow(c *config.Config, workflow config.YAMLWorkflow) ([]jobs.JobResult, error) {
	j, err := GetJobsForWorkflow(c, workflow)
	if err != nil {
		return nil, err
	}

	provider := jobs.NewInMemoryJobProvider(j)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	a := agent.NewAgent(provider, &jobs.DockerJobExecutor{}, stop, &wg)

	err = a.Start()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)
	wg.Wait()
	close(stop)
	time.Sleep(time.Second)

	return provider.Results(), nil
}
