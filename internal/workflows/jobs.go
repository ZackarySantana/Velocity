package workflows

import (
	"fmt"

	"github.com/zackarysantana/velocity/internal/jobs"
	"github.com/zackarysantana/velocity/internal/jobs/jobtypes"
	"github.com/zackarysantana/velocity/src/config"
)

type getJobsForWorkflowOpts struct {
	WithInstanceId *string
}

func WithInstanceId(instanceId string) func(*getJobsForWorkflowOpts) {
	return func(opts *getJobsForWorkflowOpts) {
		opts.WithInstanceId = &instanceId
	}
}

func GetJobsForWorkflow(c *config.Config, workflow config.YAMLWorkflow, opts ...func(*getJobsForWorkflowOpts)) ([]*jobs.Job, error) {
	o := &getJobsForWorkflowOpts{}

	for _, opt := range opts {
		opt(o)
	}

	j := []*jobs.Job{}
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

			var job jobs.Job
			if test.Run != nil {
				job = jobs.NewCommandJob(string(testName), *image.Image, *test.Run, nil, jobtypes.JobStatusQueued, nil)
			} else if test.Framework != nil && test.Language != nil {
				job = jobs.NewFrameworkJob(string(testName), *test.Language, *test.Framework, jobtypes.JobStatusQueued, &jobs.FrameworkJobOptions{
					Image:     image.Image,
					Directory: test.Directory,
				})
			}

			if job.GetName() == "" {
				return nil, fmt.Errorf("job failed to be created for test %s", testName)
			}

			j = append(j, &job)
		}
	}

	return j, nil
}
