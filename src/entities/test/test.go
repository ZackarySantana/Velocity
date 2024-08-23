package test

type Test[T any] struct {
	Id T `bson:"_id"`

	// These are populated when the routine is created.
	ImageId   T
	JobId     T
	RoutineId T

	Name string

	Language string
	Library  string

	Commands []Command

	Directory string
}
