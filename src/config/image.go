package config

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type ImageSection []Image

func (i *ImageSection) validateSyntax() error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for idx, image := range *i {
		catcher.Catch(image.error(idx).Wrap(image.validateSyntax()))
	}
	return catcher.Resolve()
}

func (i *ImageSection) validateIntegrity(c *Config) error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for idx, image := range *i {
		catcher.Catch(image.error(idx).Wrap(image.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (i *ImageSection) error(_ int) oops.OopsErrorBuilder {
	return oops.In("image_section")
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

func (i *Image) error(idx int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("image_name_%d", idx), i.Name).With(fmt.Sprintf("image_%d", idx), i.Image)
}
