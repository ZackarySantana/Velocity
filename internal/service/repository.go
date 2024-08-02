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
	Routine    RoutineRepository
	PutRoutine func(context.Context, *routine.Routine) error

	Job    JobRepository
	PutJob func(context.Context, *job.Job) error

	Image    ImageRepository
	PutImage func(context.Context, *image.Image) error

	Test    TestRepository
	PutTest func(context.Context, *test.Test) error
}

type RoutineRepository dataloader.Interface[string, *routine.Routine]
type JobRepository dataloader.Interface[string, *job.Job]
type ImageRepository dataloader.Interface[string, *image.Image]
type TestRepository dataloader.Interface[string, *test.Test]
