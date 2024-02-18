// Code generated from Pkl module `workflow_construct`. DO NOT EDIT.
package workflow

type Workflow struct {
	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Groups []*WorkflowGroup `pkl:"groups" bson:"groups,omitempty" json:"groups,omitempty" yaml:"groups,omitempty"`
}
