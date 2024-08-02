package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zackarysantana/velocity/src/config"
)

func TestImageSectionSyntax(t *testing.T) {
	testCases := []struct {
		name   string
		images config.ImageSection
		err    string
	}{
		{
			name: "image section with invalid image should have error",
			images: config.ImageSection{
				{
					Name: "name",
				},
			},
			err: "image is required",
		},
		{
			name:   "empty image section should have no error",
			images: config.ImageSection{},
		},
		{
			name: "image section with one image should have no error",
			images: config.ImageSection{
				{
					Name:  "name",
					Image: "image",
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			c := config.Config{Images: tC.images}

			if tC.err == "" {
				assert.NoError(t, c.Validate())
			} else {
				assert.ErrorContains(t, c.Validate(), tC.err)
			}
		})
	}
}

func TestImageSyntax(t *testing.T) {
	testCases := []struct {
		name  string
		image config.Image
		err   string
	}{
		{
			name:  "empty image should have no name error",
			image: config.Image{},
			err:   "name is required",
		},
		{
			name:  "empty image should have no image error",
			image: config.Image{},
			err:   "image is required",
		},
		{
			name: "image with only image should have no name error",
			image: config.Image{
				Image: "image",
			},
			err: "name is required",
		},
		{
			name: "image with only name should have no image error",
			image: config.Image{
				Name: "name",
			},
			err: "image is required",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			c := config.Config{Images: config.ImageSection{tC.image}}

			if tC.err == "" {
				assert.NoError(t, c.Validate())
			} else {
				assert.ErrorContains(t, c.Validate(), tC.err)
			}
		})
	}
}
