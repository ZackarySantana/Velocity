package image

type Image struct {
	Id string `bson:"_id"`

	Name string

	Image string
}
