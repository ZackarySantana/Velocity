// Code generated from Pkl module `build_construct`. DO NOT EDIT.
package build

import "github.com/zackarysantana/velocity/gen/pkl/primitives/command"

type Build struct {
	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	BuildRuntime string `pkl:"build_runtime" bson:"build_runtime,omitempty" json:"build_runtime,omitempty" yaml:"build_runtime,omitempty"`

	Output string `pkl:"output" bson:"output,omitempty" json:"output,omitempty" yaml:"output,omitempty"`

	OutputRuntime *string `pkl:"output_runtime" bson:"output_runtime,omitempty" json:"output_runtime,omitempty" yaml:"output_runtime,omitempty"`

	OutputCommand *string `pkl:"output_command" bson:"output_command,omitempty" json:"output_command,omitempty" yaml:"output_command,omitempty"`

	Commands []command.Command `pkl:"commands" bson:"commands,omitempty" json:"commands,omitempty" yaml:"commands,omitempty"`
}
