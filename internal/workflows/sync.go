package workflows

import (
	"sync"
	"time"

	"github.com/zackarysantana/velocity/internal/agent"
	"github.com/zackarysantana/velocity/internal/jobs"
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
				j = append(j, &jobs.CommandJob{
					Image:   *image.Image,
					Command: *test.Run,
					Name:    string(testName),
				})

				continue
			}

			j = append(j, &jobs.FrameworkJob{
				Language:  *test.Language,
				Framework: *test.Framework,
				Image:     image.Image,
				Name:      string(testName),
			})
		}
	}

	provider := jobs.NewInMemoryJobProvider(j)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	ctx := jobs.NewContext("https://github.com/zackarysantana/velocity.git", "c8dc99dfc0b62842b0a524fe34112c3df27f7e86")
	a := agent.NewAgent(provider, &jobs.DockerJobExecutor{}, ctx, stop, &wg)

	err := a.Start()
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)
	wg.Wait()
	close(stop)
	time.Sleep(time.Second)

	return provider.Results(), nil
}
