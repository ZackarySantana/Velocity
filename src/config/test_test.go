package config_test

import (
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/src/config"
)

func TestTestSectionSyntax(t *testing.T) {
	t.Run("section propogates error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.TestSection{
			{
				Name:    "name",
				Library: "library", // without language this should error
			},
			{
				Name:     "name-2",
				Language: "language",
				Commands: []config.Command{
					{
						Shell: "shell",
					},
				},
			},
			{
				Name: "name-3",
				Commands: []config.Command{
					{
						Shell: "shell",
						Params: map[string]interface{}{
							"key": "value",
						},
					},
				},
			},
		})
		require.Error(t, err)
		assert.EqualError(t, err, "validating syntax: [index=0]: language or commands are required\n[index=1]: cannot specify both a language and commands\n[index=2]: cannot specify params without a prebuilt command")

		t.Run("error fields propogate", func(t *testing.T) {
			oops, ok := oops.AsOops(err)
			require.True(t, ok)
			assert.Equal(t, "test_section", oops.Domain())

			// Test 1
			assert.Equal(t, "name", oops.Context()["test_name_0"])
			assert.Equal(t, "", oops.Context()["language_0"])
			assert.Equal(t, "library", oops.Context()["library_0"])
			assert.Empty(t, oops.Context()["commands_0"])

			// Test 2
			assert.Equal(t, "name-2", oops.Context()["test_name_1"])
			assert.Equal(t, "language", oops.Context()["language_1"])
			assert.Equal(t, "", oops.Context()["library_1"])
			assert.Len(t, oops.Context()["commands_1"].([]config.Command), 1)

			// Test 3
			assert.Equal(t, "name-3", oops.Context()["test_name_2"])
			assert.Equal(t, "", oops.Context()["language_2"])
			assert.Equal(t, "", oops.Context()["library_2"])
			assert.Len(t, oops.Context()["commands_2"].([]config.Command), 1)
		})
	})

	t.Run("section with no tests should have no error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.TestSection{})
		require.NoError(t, err)
	})
}

func TestTestSyntax(t *testing.T) {
	t.Run("empty test should have no name error and language or commands are required error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Test{})
		require.NotNil(t, err)
	})
}
