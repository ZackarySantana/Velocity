package config

import (
	"fmt"

	"github.com/samber/oops"
	"github.com/zackarysantana/velocity/src/catcher"
)

type ImageSection []Image

func (i *ImageSection) validateSyntax() error {
	return ValidateSyntaxMany(toValidators(i))
}

func (i *ImageSection) validateIntegrity(c *Config) error {
	if i == nil {
		return oops.Errorf("image section must exist and contain at least one image")
	}
	catcher := catcher.New()
	catcher.When(len(*i) == 0, "at least one image is required")
	names := make(map[string]int)
	for idx, image := range *i {
		idx2, ok := names[image.Name]
		if ok {
			catcher.Wrap(oops.Errorf("duplicate image name: %s", image.Name), "[index=%d, index_2=%d]", idx, idx2)
		}
		names[image.Name] = idx
	}
	catcher.Catch(ValidateIntegrityMany(toValidators(i), c))
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
	catcher.When(i.Name == "", "name is required")
	catcher.When(i.Image == "", "image is required")
	return catcher.Resolve()
}

func (i *Image) validateIntegrity(config *Config) error {
	return nil
}

func (i *Image) error(idx int) oops.OopsErrorBuilder {
	return oops.With(fmt.Sprintf("image_name_%d", idx), i.Name).With(fmt.Sprintf("image_%d", idx), i.Image)
}
