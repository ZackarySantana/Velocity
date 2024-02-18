// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type ShellCommand interface {
	Command

	GetCommand() string
}

var _ ShellCommand = (*ShellCommandImpl)(nil)

type ShellCommandImpl struct {
	Command string `pkl:"command" bson:"command,omitempty" json:"command,omitempty" yaml:"command,omitempty"`

	WorkingDirectory *string `pkl:"working_directory" bson:"working_directory,omitempty" json:"working_directory,omitempty" yaml:"working_directory,omitempty"`

	Env *map[string]string `pkl:"env" bson:"env,omitempty" json:"env,omitempty" yaml:"env,omitempty"`
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
