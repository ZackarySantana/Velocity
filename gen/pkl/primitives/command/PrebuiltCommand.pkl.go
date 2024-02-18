// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type PrebuiltCommand interface {
	Command

	GetPrebuilt() string

	GetParams() *map[string]any
}

var _ PrebuiltCommand = (*PrebuiltCommandImpl)(nil)

type PrebuiltCommandImpl struct {
	Prebuilt string `pkl:"prebuilt"`

	Params *map[string]any `pkl:"params"`

	WorkingDirectory *string `pkl:"working_directory"`

	Env *map[string]string `pkl:"env"`
}

func (rcv *PrebuiltCommandImpl) GetPrebuilt() string {
	return rcv.Prebuilt
}

func (rcv *PrebuiltCommandImpl) GetParams() *map[string]any {
	return rcv.Params
}

func (rcv *PrebuiltCommandImpl) GetWorkingDirectory() *string {
	return rcv.WorkingDirectory
}

func (rcv *PrebuiltCommandImpl) GetEnv() *map[string]string {
	return rcv.Env
}
