// Code generated from Pkl module `deployment_construct`. DO NOT EDIT.
package deployment

import "github.com/zackarysantana/velocity/gen/pkl/primitives/command"

type Deployment struct {
	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Runtime string `pkl:"runtime" bson:"runtime,omitempty" json:"runtime,omitempty" yaml:"runtime,omitempty"`

	Workflows *[]string `pkl:"workflows" bson:"workflows,omitempty" json:"workflows,omitempty" yaml:"workflows,omitempty"`

	Commands []command.Command `pkl:"commands" bson:"commands,omitempty" json:"commands,omitempty" yaml:"commands,omitempty"`
}
