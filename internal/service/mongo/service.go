package domain

import (
	"fmt"

	"github.com/zackarysantana/velocity/internal/service"
)

type Service struct {
	db *service.Repository
}

func NewService(db *service.Repository) (service.Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) StartRoutine() {
	// r, err := s.db.Routine.Load(context.TODO(), "test").Unwrap()

	fmt.Println("Testing this")
}
