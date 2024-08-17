package config

import (
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type CreateEntityOptions struct {
	Ic service.IdCreator

	FilterToRoutine string
}

func (c *Config) CreateEntity(opts CreateEntityOptions) (*entities.ConfigEntity, error) {
	ec := &entities.ConfigEntity{}

	for _, cRoutine := range c.Routines {
		if opts.FilterToRoutine != "" && cRoutine.Name != opts.FilterToRoutine {
			continue
		}
		other, err := c.createConfigEntityForRoutine(cRoutine, opts)
		if err != nil {
			return nil, err
		}
		ec.Merge(other)
	}

	return ec, nil
}

func (c *Config) createConfigEntityForRoutine(cRoutine Routine, opts CreateEntityOptions) (*entities.ConfigEntity, error) {
	ec := &entities.ConfigEntity{}
	eRoutine := routine.Routine{
		Id:   opts.Ic(),
		Name: cRoutine.Name,
		Jobs: []string{},
	}

	for _, jobName := range cRoutine.Jobs {
		cJob, err := c.GetJob(jobName)
		if err != nil {
			return nil, err
		}
		eJob := job.Job{
			Id:     opts.Ic(),
			Name:   cJob.Name,
			Tests:  []string{},
			Images: []string{},
		}

		for _, imageName := range cJob.Images {
			cImage, err := c.GetImage(imageName)
			if err != nil {
				return nil, err
			}
			eImage := image.Image{
				Id:    opts.Ic(),
				Name:  cImage.Name,
				Image: cImage.Image,
			}

			for _, testName := range cJob.Tests {
				cTest, err := c.GetTest(testName)
				if err != nil {
					return nil, err
				}
				eCommands := make([]test.Command, len(cTest.Commands))
				for i, cCommand := range cTest.Commands {
					eCommands[i] = test.Command{
						Shell:    cCommand.Shell,
						Prebuilt: cCommand.Prebuilt,
						Params:   cCommand.Params,
					}
				}
				eTest := test.Test{
					Id:        opts.Ic(),
					RoutineId: eRoutine.Id,
					JobId:     eJob.Id,
					ImageId:   eImage.Id,

					Name:      cTest.Name,
					Language:  cTest.Language,
					Library:   cTest.Library,
					Commands:  eCommands,
					Directory: cTest.Directory,
				}

				eJob.Tests = append(eJob.Tests, eTest.Id)
				ec.Tests = append(ec.Tests, &eTest)
			}

			eJob.Images = append(eJob.Images, eImage.Id)
			ec.Images = append(ec.Images, &eImage)
		}

		eRoutine.Jobs = append(eRoutine.Jobs, eJob.Id)
		ec.Jobs = append(ec.Jobs, &eJob)
	}

	ec.Routines = append(ec.Routines, &eRoutine)
	return ec, nil
}
