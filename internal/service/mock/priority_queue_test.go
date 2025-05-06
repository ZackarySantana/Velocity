package mock

import (
	"testing"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/servicetest"
)

func TestMongoPriorityQueue(t *testing.T) {
	pqGen := func() service.PriorityQueue[any, string] {
		return NewPriorityQueue[any, string](NewIDCreator[any]())
	}

	servicetest.TestPriorityQueue(t, pqGen)
}
