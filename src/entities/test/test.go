package test

// TODO: Add the other ids (routine, job, image)
// It will be faster to update / look them up this way.
type Test struct {
	Id string `bson:"_id"`

	Name string

	Language string
	Library  string

	Commands []Command

	Directory string
}
