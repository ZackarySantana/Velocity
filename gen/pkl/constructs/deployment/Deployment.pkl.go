// Code generated from Pkl module `deployment_construct`. DO NOT EDIT.
package deployment

import "github.com/zackarysantana/velocity/gen/pkl/primitives/command"

type Deployment struct {
	Name string `pkl:"name"`

	Runtime string `pkl:"runtime"`

	Workflows *[]string `pkl:"workflows"`

	Commands []command.Command `pkl:"commands"`
}
