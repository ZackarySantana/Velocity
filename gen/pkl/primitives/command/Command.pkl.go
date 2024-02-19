// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type Command interface {
	GetWorkingDirectory() *string

	GetEnv() *map[string]string
}
