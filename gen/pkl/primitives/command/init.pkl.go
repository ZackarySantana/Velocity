// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("command_primitive", CommandPrimitive{})
	pkl.RegisterMapping("command_primitive#PrebuiltCommand", PrebuiltCommandImpl{})
	pkl.RegisterMapping("command_primitive#ShellCommand", ShellCommandImpl{})
	pkl.RegisterMapping("command_primitive#OperationCommand", OperationCommandImpl{})
}
