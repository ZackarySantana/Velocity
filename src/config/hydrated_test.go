package config_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zackarysantana/velocity/internal/ptr"
	"github.com/zackarysantana/velocity/src/config"
)

func TestHydrateWorkflowGroup(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawWorkflowGroup
		hydrated config.WorkflowGroup
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawWorkflowGroup{},
			hydrated: config.WorkflowGroup{},
		},
		{
			name: "with name",
			raw: config.RawWorkflowGroup{
				Name: "node group",
			},
			hydrated: config.WorkflowGroup{
				Name: "node group",
			},
		},
		{
			name: "with runtimes",
			raw: config.RawWorkflowGroup{
				Name:     "node group",
				Runtimes: []string{"node"},
			},
			hydrated: config.WorkflowGroup{
				Name:     "node group",
				Runtimes: []string{"node"},
			},
		},
		{
			name: "with tests",
			raw: config.RawWorkflowGroup{
				Name:  "node group",
				Tests: []string{"test"},
			},
			hydrated: config.WorkflowGroup{
				Name:  "node group",
				Tests: []string{"test"},
			},
		},
		{
			name: "with runtimes and tests",
			raw: config.RawWorkflowGroup{
				Name:     "node group",
				Runtimes: []string{"node"},
				Tests:    []string{"test"},
			},
			hydrated: config.WorkflowGroup{
				Name:     "node group",
				Runtimes: []string{"node"},
				Tests:    []string{"test"},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateWorkflowGroup(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateWorkflow(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawWorkflow
		hydrated config.Workflow
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawWorkflow{},
			hydrated: config.Workflow{},
		},
		{
			name: "with name",
			raw: config.RawWorkflow{
				Name: "workflow",
			},
			hydrated: config.Workflow{
				Name: "workflow",
			},
		},
		{
			name: "with group",
			raw: config.RawWorkflow{
				Name: "workflow",
				Groups: []config.RawWorkflowGroup{
					{
						Name:     "node group",
						Runtimes: []string{"node"},
						Tests:    []string{"test"},
					},
				},
			},
			hydrated: config.Workflow{
				Name: "workflow",
				Groups: []config.WorkflowGroup{
					{
						Name:     "node group",
						Runtimes: []string{"node"},
						Tests:    []string{"test"},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateWorkflow(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}

}

func TestHydrateDeployment(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawDeployment
		hydrated config.Deployment
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawDeployment{},
			hydrated: config.Deployment{},
		},
		{
			name: "with required arguments",
			raw: config.RawDeployment{
				Name: "deploy",
			},
			hydrated: config.Deployment{
				Name: "deploy",
			},
		},
		{
			name: "with workflows",
			raw: config.RawDeployment{
				Name:      "deploy",
				Workflows: []string{"workflow"},
			},
			hydrated: config.Deployment{
				Name:      "deploy",
				Workflows: []string{"workflow"},
			},
		},
		{
			name: "with commands",
			raw: config.RawDeployment{
				Name: "deploy",
				Commands: []config.RawCommand{
					{
						Command: ptr.To("echo 'hello world'"),
					},
					{
						Command: ptr.To("echo '2nd command'"),
					},
				},
			},
			hydrated: config.Deployment{
				Name: "deploy",
				Commands: []config.Command{
					config.ShellCommand{
						Command: "echo 'hello world'",
					},
					config.ShellCommand{
						Command: "echo '2nd command'",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateDeployment(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateBuild(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawBuild
		hydrated config.Build
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawBuild{},
			hydrated: config.Build{},
		},
		{
			name: "with required arguments",
			raw: config.RawBuild{
				Name:         "build",
				BuildRuntime: "node",
				Output:       "dist",
			},
			hydrated: config.Build{
				Name:         "build",
				BuildRuntime: "node",
				Output:       "dist",
			},
		},
		{
			name: "with output arguments",
			raw: config.RawBuild{
				Name:          "build",
				BuildRuntime:  "node",
				Output:        "dist",
				OutputRuntime: ptr.To("node2"),
				OutputCmd:     ptr.To("npm start"),
			},
			hydrated: config.Build{
				Name:          "build",
				BuildRuntime:  "node",
				Output:        "dist",
				OutputRuntime: ptr.To("node2"),
				OutputCmd:     ptr.To("npm start"),
			},
		},
		{
			name: "with commands",
			raw: config.RawBuild{
				Name:         "build",
				BuildRuntime: "node",
				Output:       "dist",
				Commands: []config.RawCommand{
					{
						Command: ptr.To("echo 'hello world'"),
					},
					{
						Command: ptr.To("echo '2nd command'"),
					},
				},
			},
			hydrated: config.Build{
				Name:         "build",
				BuildRuntime: "node",
				Output:       "dist",
				Commands: []config.Command{
					config.ShellCommand{
						Command: "echo 'hello world'",
					},
					config.ShellCommand{
						Command: "echo '2nd command'",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateBuild(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateRuntime(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawRuntime
		hydrated config.Runtime
		err      error
	}{
		{
			name: "empty",
			raw:  config.RawRuntime{},
			err:  fmt.Errorf("invalid runtime: %v", config.RawRuntime{}),
		},
		{
			name: "docker",
			raw: config.RawRuntime{
				Name:  "docker",
				Image: ptr.To("node:latest"),
			},
			hydrated: config.DockerRuntime{
				Name_: "docker",
				Image: "node:latest",
			},
		},
		{
			name: "machine",
			raw: config.RawRuntime{
				Name:    "machine",
				Machine: ptr.To("linux"),
			},
			hydrated: config.BareMetalRuntime{
				Name_:   "machine",
				Machine: ptr.To("linux"),
			},
		},
		{
			name: "with env",
			raw: config.RawRuntime{
				Name:  "docker",
				Image: ptr.To("node:latest"),
				Env:   &config.RawEnv{"APP=app"},
			},
			hydrated: config.DockerRuntime{
				Name_: "docker",
				Image: "node:latest",
				Env_: &config.Env{
					"APP": "app",
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateRuntime(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateTest(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawTest
		hydrated config.Test
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawTest{},
			hydrated: config.Test{},
		},
		{
			name: "only name",
			raw: config.RawTest{
				Name: "test",
			},
			hydrated: config.Test{
				Name: "test",
			},
		},
		{
			name: "with env",
			raw: config.RawTest{
				Name: "test",
				Env:  &config.RawEnv{"APP=app"},
			},
			hydrated: config.Test{
				Name: "test",
				Env:  &config.Env{"APP": "app"},
			},
		},
		{
			name: "with wd",
			raw: config.RawTest{
				Name:             "test",
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.Test{
				Name:             "test",
				WorkingDirectory: ptr.To("/app"),
			},
		},
		{
			name: "with env and wd",
			raw: config.RawTest{
				Name:             "test",
				Env:              &config.RawEnv{"APP=app"},
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.Test{
				Name:             "test",
				Env:              &config.Env{"APP": "app"},
				WorkingDirectory: ptr.To("/app"),
			},
		},
		{
			name: "with commands",
			raw: config.RawTest{
				Name: "test",
				Commands: []config.RawCommand{
					{
						Command: ptr.To("echo 'hello world'"),
					},
					{
						Command: ptr.To("echo '2nd command'"),
					},
				},
			},
			hydrated: config.Test{
				Name: "test",
				Commands: []config.Command{
					config.ShellCommand{
						Command: "echo 'hello world'",
					},
					config.ShellCommand{
						Command: "echo '2nd command'",
					},
				},
			},
		},
		{
			name: "with env, wd, and commands",
			raw: config.RawTest{
				Name:             "test",
				Env:              &config.RawEnv{"APP=app"},
				WorkingDirectory: ptr.To("/app"),
				Commands: []config.RawCommand{
					{
						Command: ptr.To("echo 'hello world'"),
					},
					{
						Command: ptr.To("echo '2nd command'"),
					},
				},
			},
			hydrated: config.Test{
				Name:             "test",
				Env:              &config.Env{"APP": "app"},
				WorkingDirectory: ptr.To("/app"),
				Commands: []config.Command{
					config.ShellCommand{
						Command: "echo 'hello world'",
					},
					config.ShellCommand{
						Command: "echo '2nd command'",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateTest(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateOperationCommand(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawCommand
		hydrated config.OperationCommand
		err      error
	}{
		{
			name: "empty",
			raw:  config.RawCommand{},
			err:  fmt.Errorf("invalid command: %v", config.RawCommand{}),
		},
		{
			name: "only command",
			raw: config.RawCommand{
				Operation: ptr.To("echo 'hello world'"),
			},
			hydrated: config.OperationCommand{
				Operation: "echo 'hello world'",
			},
		},
		{
			name: "with env",
			raw: config.RawCommand{
				Operation: ptr.To("echo 'hello world'"),
				Env:       &config.RawEnv{"APP=app"},
			},
			hydrated: config.OperationCommand{
				Operation: "echo 'hello world'",
				Env_:      &config.Env{"APP": "app"},
			},
		},
		{
			name: "with wd",
			raw: config.RawCommand{
				Operation:        ptr.To("echo 'hello world'"),
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.OperationCommand{
				Operation:         "echo 'hello world'",
				WorkingDirectory_: ptr.To("/app"),
			},
		},
		{
			name: "with env and wd",
			raw: config.RawCommand{
				Operation:        ptr.To("echo 'hello world'"),
				Env:              &config.RawEnv{"APP=app"},
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.OperationCommand{
				Operation:         "echo 'hello world'",
				Env_:              &config.Env{"APP": "app"},
				WorkingDirectory_: ptr.To("/app"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateOperationCommand(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateShellCommand(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawCommand
		hydrated config.ShellCommand
		err      error
	}{
		{
			name: "empty",
			raw:  config.RawCommand{},
			err:  fmt.Errorf("invalid command: %v", config.RawCommand{}),
		},
		{
			name: "only command",
			raw: config.RawCommand{
				Command: ptr.To("echo 'hello world'"),
			},
			hydrated: config.ShellCommand{
				Command: "echo 'hello world'",
			},
		},
		{
			name: "with env",
			raw: config.RawCommand{
				Command: ptr.To("echo 'hello world'"),
				Env:     &config.RawEnv{"APP=app"},
			},
			hydrated: config.ShellCommand{
				Command: "echo 'hello world'",
				Env_:    &config.Env{"APP": "app"},
			},
		},
		{
			name: "with wd",
			raw: config.RawCommand{
				Command:          ptr.To("echo 'hello world'"),
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.ShellCommand{
				Command:           "echo 'hello world'",
				WorkingDirectory_: ptr.To("/app"),
			},
		},
		{
			name: "with env and wd",
			raw: config.RawCommand{
				Command:          ptr.To("echo 'hello world'"),
				Env:              &config.RawEnv{"APP=app"},
				WorkingDirectory: ptr.To("/app"),
			},
			hydrated: config.ShellCommand{
				Command:           "echo 'hello world'",
				Env_:              &config.Env{"APP": "app"},
				WorkingDirectory_: ptr.To("/app"),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateShellCommand(tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			assert.Equal(tt.hydrated, hydrated)
		})
	}
}

func TestHydrateEnv(t *testing.T) {
	tests := []struct {
		name     string
		raw      config.RawEnv
		hydrated config.Env
		err      error
	}{
		{
			name:     "empty",
			raw:      config.RawEnv{},
			hydrated: config.Env{},
		},
		{
			name: "with env",
			raw: config.RawEnv{
				"APP=app",
			},
			hydrated: config.Env{
				"APP": "app",
			},
		},
		{
			name: "with spaces",
			raw: config.RawEnv{
				"APP= app app",
			},
			hydrated: config.Env{
				"APP": " app app",
			},
		},
		{
			name: "with escaped quotes",
			raw: config.RawEnv{
				"APP=\\\"app\\\"",
			},
			hydrated: config.Env{
				"APP": "\\\"app\\\"",
			},
		},
		{
			name: "with escaped quotes surrounded by quotes",
			raw: config.RawEnv{
				"APP=\"\\\"app\\\"\"",
			},
			hydrated: config.Env{
				"APP": "\\\"app\\\"",
			},
		},
		{
			name: "with multiple env",
			raw: config.RawEnv{
				"APP=app",
				"APP2=app2",
			},
			hydrated: config.Env{
				"APP":  "app",
				"APP2": "app2",
			},
		},
		{
			name: "with single quotes",
			raw: config.RawEnv{
				"APP='app'",
			},
			hydrated: config.Env{
				"APP": "app",
			},
		},
		{
			name: "with double quotes",
			raw: config.RawEnv{
				"APP=\"app\"",
			},
			hydrated: config.Env{
				"APP": "app",
			},
		},
		{
			name: "with no value and equal",
			raw: config.RawEnv{
				"APP",
			},
			err: errors.New("invalid env line: APP"),
		},
		{
			name: "with no value",
			raw: config.RawEnv{
				"APP=",
			},
			hydrated: config.Env{
				"APP": "",
			},
		},
		{
			name: "with no name",
			raw: config.RawEnv{
				"=app",
			},
			err: errors.New("invalid env line: =app"),
		},
		{
			name: "with imbalanced quotes",
			raw: config.RawEnv{
				"APP='app\"",
			},
			err: errors.New("invalid env line: APP='app\""),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert := assert.New(t)
			require := require.New(t)
			hydrated, err := config.HydrateEnv(&tt.raw)
			if tt.err != nil {
				require.EqualError(err, tt.err.Error())
				return
			} else {
				require.NoError(err)
			}
			require.NotNil(hydrated)
			assert.Equal(tt.hydrated, *hydrated)
		})
	}
}
