// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type ShellCommand interface {
	Command

	GetCommand() string
}

var _ ShellCommand = (*ShellCommandImpl)(nil)

type ShellCommandImpl struct {
	Command string `pkl:"command"`

	WorkingDirectory *string `pkl:"working_directory"`

	Env *map[string]string `pkl:"env"`
}

func (rcv *ShellCommandImpl) GetCommand() string {
	return rcv.Command
}

func (rcv *ShellCommandImpl) GetWorkingDirectory() *string {
	return rcv.WorkingDirectory
}

func (rcv *ShellCommandImpl) GetEnv() *map[string]string {
	return rcv.Env
}
