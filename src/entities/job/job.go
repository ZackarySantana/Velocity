package job

type Job[T any] struct {
	Id T `bson:"_id"`

	Name string

	// Tests is a list of test ids.
	Tests []T

	// Images is a list of image ids.
	Images []T
}
