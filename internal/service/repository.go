package service

import (
	"context"

	"github.com/sysulq/dataloader-go"
	"github.com/zackarysantana/velocity/src/entities/image"
	"github.com/zackarysantana/velocity/src/entities/job"
	"github.com/zackarysantana/velocity/src/entities/routine"
	"github.com/zackarysantana/velocity/src/entities/test"
)

type Repository struct {
	Routine     RoutineRepository
	PutRoutines func(context.Context, []*routine.Routine) error

	Job     JobRepository
	PutJobs func(context.Context, []*job.Job) error

	Image     ImageRepository
	PutImages func(context.Context, []*image.Image) error

	Test     TestRepository
	PutTests func(context.Context, []*test.Test) error

	WithTransaction func(context.Context, func(context.Context) error) error
}

type RoutineRepository dataloader.Interface[string, *routine.Routine]
type JobRepository dataloader.Interface[string, *job.Job]
type ImageRepository dataloader.Interface[string, *image.Image]
type TestRepository dataloader.Interface[string, *test.Test]
