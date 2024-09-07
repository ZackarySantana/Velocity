package mock

import (
	"testing"

	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/internal/service/servicetest"
)

func TestMongoPriorityQueue(t *testing.T) {
	i := 0
	pqGen := func() service.PriorityQueue[any, string] {
		i++
		return NewPriorityQueue[any, string](NewIDCreator[any]())
	}

	servicetest.TestPriorityQueue(t, pqGen)
}
