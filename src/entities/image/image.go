package image

type Image[T any] struct {
	Id T `bson:"_id"`

	Name string

	Image string
}
