// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

type Runtime interface {
	GetName() string

	GetEnv() *map[string]string
}
