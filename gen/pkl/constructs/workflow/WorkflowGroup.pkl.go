// Code generated from Pkl module `workflow_construct`. DO NOT EDIT.
package workflow

type WorkflowGroup struct {
	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Runtimes []string `pkl:"runtimes" bson:"runtimes,omitempty" json:"runtimes,omitempty" yaml:"runtimes,omitempty"`

	Tests []string `pkl:"tests" bson:"tests,omitempty" json:"tests,omitempty" yaml:"tests,omitempty"`
}
