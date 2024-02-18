// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

type MachineRuntime interface {
	Runtime

	GetMachine() string
}

var _ MachineRuntime = (*MachineRuntimeImpl)(nil)

type MachineRuntimeImpl struct {
	Machine string `pkl:"machine"`

	Name string `pkl:"name"`

	Env *map[string]string `pkl:"env"`
}

func (rcv *MachineRuntimeImpl) GetMachine() string {
	return rcv.Machine
}

func (rcv *MachineRuntimeImpl) GetName() string {
	return rcv.Name
}

func (rcv *MachineRuntimeImpl) GetEnv() *map[string]string {
	return rcv.Env
}
