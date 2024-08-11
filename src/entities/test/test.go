package test

type Test struct {
	Id string `bson:"_id"`

	RoutineId string
	JobId     string
	ImageId   string

	Name string

	Language string
	Library  string

	Commands []Command

	Directory string
}
