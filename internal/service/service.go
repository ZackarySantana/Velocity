package service

import (
	"context"

	"github.com/zackarysantana/velocity/src/entities"
)

type Service interface {
	StartRoutine(context.Context, *entities.ConfigEntity, string) error
}
