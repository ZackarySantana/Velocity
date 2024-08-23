package routine

type Routine[T any] struct {
	Id T `bson:"_id"`

	Name string

	// Jobs is a list of job ids.
	Jobs []T
}
