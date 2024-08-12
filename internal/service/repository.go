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
	Routine *RoutineRepository
	Job     *JobRepository
	Image   *ImageRepository
	Test    *TestRepository

	WithTransaction func(context.Context, func(context.Context) error) error
}

type RoutineRepository struct {
	dataloader.Interface[string, *routine.Routine]
	Put func(context.Context, []*routine.Routine) error
}

type JobRepository struct {
	dataloader.Interface[string, *job.Job]
	Put func(context.Context, []*job.Job) error
}

type ImageRepository struct {
	dataloader.Interface[string, *image.Image]
	Put func(context.Context, []*image.Image) error
}

type TestRepository struct {
	dataloader.Interface[string, *test.Test]
	Put func(context.Context, []*test.Test) error
}
