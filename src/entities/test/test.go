package test

type Test struct {
	Id string `bson:"_id"`

	Name string

	Language string
	Library  string

	Commands []Command

	Directory string
}
