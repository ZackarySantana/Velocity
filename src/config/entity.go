package config

import (
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/entities"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type CreateEntityOptions[T any] struct {
	Id service.IdCreator[T]

	FilterToRoutine string
}

func CreateEntity[T any](config *Config, opts CreateEntityOptions[T]) (*entities.ConfigEntity[T], error) {
	ec := &entities.ConfigEntity[T]{}

	for _, cRoutine := range config.Routines {
		if opts.FilterToRoutine != "" && cRoutine.Name != opts.FilterToRoutine {
			continue
		}
		other, err := createConfigEntityForRoutine(config, cRoutine, opts)
		if err != nil {
			return nil, err
		}
		ec.Merge(other)
	}

	return ec, nil
}

func createConfigEntityForRoutine[T any](config *Config, cRoutine Routine, opts CreateEntityOptions[T]) (*entities.ConfigEntity[T], error) {
	ec := &entities.ConfigEntity[T]{}
	eRoutine := routine.Routine[T]{
		Id:   opts.Id.Create(),
		Name: cRoutine.Name,
		Jobs: []T{},
	}

	for _, jobName := range cRoutine.Jobs {
		cJob, err := config.GetJob(jobName)
		if err != nil {
			return nil, err
		}
		eJob := job.Job[T]{
			Id:     opts.Id.Create(),
			Name:   cJob.Name,
			Tests:  []T{},
			Images: []T{},
		}

		for _, imageName := range cJob.Images {
			cImage, err := config.GetImage(imageName)
			if err != nil {
				return nil, err
			}
			eImage := image.Image[T]{
				Id:    opts.Id.Create(),
				Name:  cImage.Name,
				Image: cImage.Image,
			}

			for _, testName := range cJob.Tests {
				cTest, err := config.GetTest(testName)
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
				eTest := test.Test[T]{
					Id:        opts.Id.Create(),
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
