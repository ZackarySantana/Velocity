package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
	"github.com/zackarysantana/velocity/src/config/id"
	"github.com/zackarysantana/velocity/src/entities/image"
)

type ImageSection []Image

func (i *ImageSection) validateSyntax() error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for _, image := range *i {
		catcher.Catch(image.error().Wrap(image.validateSyntax()))
	}
	return catcher.Resolve()
}

func (i *ImageSection) validateIntegrity(c *Config) error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for _, image := range *i {
		catcher.Catch(image.error().Wrap(image.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (i *ImageSection) error() oops.OopsErrorBuilder {
	return oops.In("image_section")
}

func (i *ImageSection) ToEntities(ic id.Creator) []*image.Image {
	images := make([]*image.Image, 0)
	for _, img := range *i {
		images = append(images, img.ToEntity(ic))
	}
	return images
}

type Image struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func (i *Image) validateSyntax() error {
	catcher := catcher.New()
	catcher.ErrorWhen(i.Name == "", "name is required")
	catcher.ErrorWhen(i.Image == "", "image is required")
	return catcher.Resolve()
}

func (i *Image) validateIntegrity(config *Config) error {
	return nil
}

func (i *Image) error() oops.OopsErrorBuilder {
	return oops.With("image_name", i.Name)
}

func (i *Image) ToEntity(ic id.Creator) *image.Image {
	return &image.Image{
		Id:    ic(),
		Name:  i.Name,
		Image: i.Image,
	}
}
