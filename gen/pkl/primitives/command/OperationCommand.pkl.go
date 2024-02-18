// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type OperationCommand interface {
	Command

	GetOperation() string
}

var _ OperationCommand = (*OperationCommandImpl)(nil)

type OperationCommandImpl struct {
	Operation string `pkl:"operation" bson:"operation,omitempty" json:"operation,omitempty" yaml:"operation,omitempty"`

	WorkingDirectory *string `pkl:"working_directory" bson:"working_directory,omitempty" json:"working_directory,omitempty" yaml:"working_directory,omitempty"`

	Env *map[string]string `pkl:"env" bson:"env,omitempty" json:"env,omitempty" yaml:"env,omitempty"`
}

func (rcv *OperationCommandImpl) GetOperation() string {
	return rcv.Operation
}

func (rcv *OperationCommandImpl) GetWorkingDirectory() *string {
	return rcv.WorkingDirectory
}

func (rcv *OperationCommandImpl) GetEnv() *map[string]string {
	return rcv.Env
}
