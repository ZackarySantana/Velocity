package routine

type Routine struct {
	Id string `bson:"_id"`

	Name string

	// Jobs is a list of job ids.
	Jobs []string
}
