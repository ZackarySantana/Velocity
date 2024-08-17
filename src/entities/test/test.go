package test

type Test struct {
	Id string `bson:"_id"`

	// These are populated when the routine is created.
	ImageId   string
	JobId     string
	RoutineId string

	Name string

	Language string
	Library  string

	Commands []Command

	Directory string
}
