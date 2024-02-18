// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

import "github.com/apple/pkl-go/pkl"

func init() {
	pkl.RegisterMapping("runtime_primitive", RuntimePrimitive{})
	pkl.RegisterMapping("runtime_primitive#DockerRuntime", DockerRuntimeImpl{})
	pkl.RegisterMapping("runtime_primitive#MachineRuntime", MachineRuntimeImpl{})
}
