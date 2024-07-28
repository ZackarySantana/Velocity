package config

import (
	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type ImageSection []Image

func (i *ImageSection) validateSyntax() error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for _, image := range *i {
		catcher.Catch(image.Error().Wrap(image.validateSyntax()))
	}
	return catcher.Resolve()
}

func (i *ImageSection) validateIntegrity(c *Config) error {
	if i == nil {
		return nil
	}
	catcher := catcher.New()
	for _, image := range *i {
		catcher.Catch(image.Error().Wrap(image.validateIntegrity(c)))
	}
	return catcher.Resolve()
}

func (i *ImageSection) Error() oops.OopsErrorBuilder {
	return oops.Code("image_section")
}

type Image struct {
	Name  string `yaml:"name"`
	Image string `yaml:"image"`
}

func (i *Image) validateSyntax() error {
	catcher := catcher.New()
	if i.Name == "" {
		catcher.Error("name is required")
	}
	if i.Image == "" {
		catcher.Error("image is required")
	}
	return catcher.Resolve()
}

func (i *Image) validateIntegrity(config *Config) error {
	return nil
}

func (i *Image) Error() oops.OopsErrorBuilder {
	return oops.With("image_name", i.Name)
}
