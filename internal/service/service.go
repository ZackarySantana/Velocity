package service

import (
	"context"

	"github.com/zackarysantana/velocity/src/entities"
)

type Service[T any] interface {
	StartRoutine(context.Context, *entities.ConfigEntity[T], string) error
}
