package config_test

import (
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/src/config"
)

func TestImageSectionSyntax(t *testing.T) {
	t.Run("section propogates error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.ImageSection{
			{
				Name: "name",
			},
			{
				Image: "image",
			},
		})
		require.Error(t, err)
		assert.EqualError(t, err, "validating syntax: [index=0]: image is required\n[index=1]: name is required")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "image_section", oops.Domain())

			// Image one
			assert.Equal(t, "name", oops.Context()["image_name_0"])
			assert.Equal(t, "", oops.Context()["image_0"])

			// Image two
			assert.Equal(t, "", oops.Context()["image_name_1"])
			assert.Equal(t, "image", oops.Context()["image_1"])
		})
	})

	t.Run("section with no images should have no error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.ImageSection{})
		require.NoError(t, err)
	})
}

func TestImageSectionIntegrity(t *testing.T) {
	t.Run("section without images should error", func(t *testing.T) {
		err := config.ValidateIntegrity(&config.ImageSection{}, &config.Config{})
		require.Error(t, err)
		assert.EqualError(t, err, "validating integrity: at least one image is required")
	})

	t.Run("section with duplicate image names should error", func(t *testing.T) {
		err := config.ValidateIntegrity(&config.ImageSection{
			{
				Name:  "name",
				Image: "image",
			},
			{
				Name:  "name",
				Image: "image_2",
			},
		}, &config.Config{})
		require.Error(t, err)
		assert.EqualError(t, err, "validating integrity: [index=1, index_2=0]: duplicate image name: name")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "image_section", oops.Domain())
		})
	})
}

func TestImageSyntax(t *testing.T) {
	t.Run("empty image should have no name error and no image error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Image{})
		require.Error(t, err)
		assert.EqualError(t, err, "validating syntax: name is required\nimage is required")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "", oops.Context()["image_name_0"])
			assert.Equal(t, "", oops.Context()["image_0"])
		})
	})

	t.Run("image with only image should have no name error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Image{
			Image: "image",
		})
		require.Error(t, err)
		assert.EqualError(t, err, "validating syntax: name is required")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "", oops.Context()["image_name_0"])
			assert.Equal(t, "image", oops.Context()["image_0"])
		})
	})

	t.Run("image with only name should have no image error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Image{
			Name: "name",
		})
		require.Error(t, err)
		assert.EqualError(t, err, "validating syntax: image is required")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "name", oops.Context()["image_name_0"])
			assert.Equal(t, "", oops.Context()["image_0"])
		})
	})

	t.Run("image with name and image should have no error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Image{
			Name:  "name",
			Image: "image",
		})
		require.NoError(t, err)
	})
}
