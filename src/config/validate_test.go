package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/src/config"
)

var (
	validConfig = `
config:
    registry: docker.io

tests:
    t-test-1:
        directory: tests/end-to-end
        exclude_directories:
            - tests
        language: golang
        framework: std

    t-test-2:
        directory: tests/end-to-end
        run: echo "TBA"

images:
    ubuntu:
        image: ubuntu:mantic

workflows:
    t-workflow-1:
        tests:
            ubuntu:
                - t-test-1
                - t-test-2

    t-workflow-2:
        tests:
            ubuntu:
                - t-test-2
`
	// Test is missing framework
	invalidTest1Config = `
config:
    registry: docker.io

tests:
    t-test-1:
        directory: tests/end-to-end
        exclude_directories:
            - tests
        language: golang

images:
    ubuntu:
        image: ubuntu:mantic

workflows:
    t-workflow-1:
        tests:
            ubuntu:
                - t-test-1
`
	// Test has both Run and Language
	invalidTest2Config = `
config:
    registry: docker.io

tests:
    t-test-1:
        directory: tests/end-to-end
        exclude_directories:
            - tests
        language: golang
        run: echo "TEST"

images:
    ubuntu:
        image: ubuntu:mantic

workflows:
    t-workflow-1:
        tests:
            ubuntu:
                - t-test-1
`
)

func TestMain(m *testing.M) {
	os.Setenv("VELOCITY_CONFIG", "../../velocity.yml")
	m.Run()
}

func TestLoadConfig(t *testing.T) {
	// Test local file
	_, err := config.LoadConfig()
	assert.NoError(t, err)

	// Test remote file
	os.Setenv("VELOCITY_CONFIG", "https://raw.githubusercontent.com/zackarysantana/velocity/main/velocity.yml")
	_, err = config.LoadConfig()
	assert.NoError(t, err)
}

func TestParseConfig(t *testing.T) {
	parser := config.MultiParser{&config.YAMLParser{}}

	// Parsing valid config
	_, err := config.ParseConfig([]byte(validConfig), parser)
	assert.NoError(t, err)

	// Parsing invalid config that does not have language and framework
	_, err = config.ParseConfig([]byte(invalidTest1Config), parser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test 't-test-1' must have either language & framework or run")

	// Parsing invalid config that does has language and run
	_, err = config.ParseConfig([]byte(invalidTest2Config), parser)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "test 't-test-1' must have either language & framework or run- not both")
}

func TestGetWorkflow(t *testing.T) {
	c, err := config.LoadConfig()
	require.NoError(t, err)
	_, err = c.GetWorkflow("velocity")
	assert.NoError(t, err)
}

func TestGetImage(t *testing.T) {
	c, err := config.LoadConfig()
	require.NoError(t, err)

	_, err = c.GetImage("ubuntu")
	assert.NoError(t, err)
}

func TestGetTest(t *testing.T) {
	c, err := config.LoadConfig()
	require.NoError(t, err)

	_, err = c.GetTest("velocity")
	assert.NoError(t, err)
}

func TestValidate(t *testing.T) {
	c, err := config.LoadConfig()
	require.NoError(t, err)

	err = c.Validate()
	assert.NoError(t, err)
}
