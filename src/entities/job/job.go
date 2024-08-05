package job

type Job struct {
	Id string `bson:"_id"`

	Name string

	// Tests is a list of test ids.
	Tests []string

	// Images is a list of image ids.
	Images []string
}
