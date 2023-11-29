package workflows

import (
	"os/exec"
	"strings"
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
					Image:     *image.Image,
					Directory: test.Directory,
					Command:   *test.Run,
					Name:      string(testName),
				})

				continue
			}

			j = append(j, &jobs.FrameworkJob{
				Image:     image.Image,
				Directory: test.Directory,
				Language:  *test.Language,
				Framework: *test.Framework,
				Name:      string(testName),
			})
		}
	}

	gitRepo, err := exec.Command("git", "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return nil, err
	}

	commitHash, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return nil, err
	}

	provider := jobs.NewInMemoryJobProvider(j)

	stop := make(chan bool)
	wg := sync.WaitGroup{}
	repo, _ := strings.CutSuffix(string(gitRepo), "\n")
	hash, _ := strings.CutSuffix(string(commitHash), "\n")
	ctx := jobs.NewContext(repo, hash)
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
