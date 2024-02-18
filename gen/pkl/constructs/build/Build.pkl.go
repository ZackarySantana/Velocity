// Code generated from Pkl module `build_construct`. DO NOT EDIT.
package build

import "github.com/zackarysantana/velocity/gen/pkl/primitives/command"

type Build struct {
	Name string `pkl:"name"`

	BuildRuntime string `pkl:"build_runtime"`

	Output string `pkl:"output"`

	OutputRuntime *string `pkl:"output_runtime"`

	OutputCommand *string `pkl:"output_command"`

	Commands []command.Command `pkl:"commands"`
}
