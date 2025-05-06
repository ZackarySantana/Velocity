package mock

import (
	"testing"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/servicetest"
)

func TestMockRepository(t *testing.T) {
	repoGen := func() service.RepositoryManager[any] {
		return NewRepositoryManager(NewIDCreator[any]())
	}

	servicetest.TestRepository(t, repoGen)
}
