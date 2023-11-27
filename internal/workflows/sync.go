package workflows

import (
	"fmt"
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/src/config"
)

func RunSyncWorkflow(c config.Config, workflow config.YAMLWorkflow) ([]jobs.JobResult, error) {
	j := []jobs.Job{}
	for image, testNames := range workflow.Tests {
		for _, testName := range testNames {
			test, err := c.GetTest(string(testName))
			if err != nil {
				return nil, err
			}
			// TODO: How do we run everything possible

			run := "echo 'hello world'"

			if test.Run != nil {
				run = *test.Run
			}

			fmt.Println("Running test '" + test.Name + "' on image '" + image + "'")

			j = append(j, jobs.Job{
				Image:   image,
				Command: run,
			})
		}
	}

	provider := jobs.NewInMemoryJobProvider(j)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	a := agent.NewAgent(provider, &jobs.DockerJobExecutor{}, stop, &wg)

	err := a.Start()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	wg.Wait()
	close(stop)

	return provider.Results(), nil
}
