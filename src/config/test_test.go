package config_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/src/config"
)

func TestTestSectionSyntax(t *testing.T) {

}

func TestTestSyntax(t *testing.T) {
	t.Run("empty test should have no name error and language or commands are required error", func(t *testing.T) {
		err := config.ValidateSyntax(&config.Test{})
		require.NotNil(t, err)
	})
}
