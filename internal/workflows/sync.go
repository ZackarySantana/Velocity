package workflows

import (
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"github.com/zackarysantana/velocity/src/config"
)

func RunSyncWorkflow(c config.Config, workflow config.YAMLWorkflow) ([]jobs.JobResult, error) {
	j := []jobs.Job{}
	for imageName, testNames := range workflow.Tests {
		for _, testName := range testNames {
			test, err := c.GetTest(string(testName))
			if err != nil {
				return nil, err
			}
			image, err := c.GetImage(imageName)
			if err != nil {
				return nil, err
			}

			if test.Run != nil {
				j = append(j, jobs.NewCommandJob(string(testName), *image.Image, *test.Run, nil, jobtypes.JobStatusQueued, nil))

				continue
			}

			j = append(j, jobs.NewFrameworkJob(string(testName), *test.Language, *test.Framework, jobtypes.JobStatusQueued, &jobs.FrameworkJobOptions{
				Image:     image.Image,
				Directory: test.Directory,
			}))
		}
	}

	provider := jobs.NewInMemoryJobProvider(j)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	ctx, err := jobs.NewCurrentContext()
	if err != nil {
		return nil, err
	}
	a := agent.NewAgent(provider, &jobs.DockerJobExecutor{}, ctx, stop, &wg)

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
