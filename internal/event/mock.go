package event

import "context"

type Mock struct{}

func NewMock() EventSender {
	return &Mock{}
}

func (m *Mock) SendEvent(ctx context.Context, event Event) error {
	return nil
}
