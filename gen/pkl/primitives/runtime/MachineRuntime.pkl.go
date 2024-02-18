// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

type MachineRuntime interface {
	Runtime

	GetMachine() string
}

var _ MachineRuntime = (*MachineRuntimeImpl)(nil)

type MachineRuntimeImpl struct {
	Machine string `pkl:"machine" bson:"machine,omitempty" json:"machine,omitempty" yaml:"machine,omitempty"`

	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Env *map[string]string `pkl:"env" bson:"env,omitempty" json:"env,omitempty" yaml:"env,omitempty"`
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
