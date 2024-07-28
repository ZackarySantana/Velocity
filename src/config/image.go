package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type ImageSection []Image

func (i *ImageSection) Validate() error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for _, image := range *i {
		catcher.Catch(image.Validate())
	}
	return catcher.Resolve()
}

type Image struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func (i *Image) Validate() error {
	oops := oops.With("image", *i)
	if i.Name == "" {
		return oops.Errorf("name is required")
	}
	if i.Image == "" {
		return oops.Errorf("image is required")
	}
	return nil
}
