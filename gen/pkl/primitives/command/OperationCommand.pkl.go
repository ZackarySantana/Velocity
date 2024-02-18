// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type OperationCommand interface {
	Command

	GetOperation() string
}

var _ OperationCommand = (*OperationCommandImpl)(nil)

type OperationCommandImpl struct {
	Operation string `pkl:"operation"`

	WorkingDirectory *string `pkl:"working_directory"`

	Env *map[string]string `pkl:"env"`
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
