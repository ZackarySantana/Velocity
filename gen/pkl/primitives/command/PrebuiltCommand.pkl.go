// Code generated from Pkl module `command_primitive`. DO NOT EDIT.
package command

type PrebuiltCommand interface {
	Command

	GetPrebuilt() string

	GetParams() *map[string]any
}

var _ PrebuiltCommand = (*PrebuiltCommandImpl)(nil)

type PrebuiltCommandImpl struct {
	Prebuilt string `pkl:"prebuilt" bson:"prebuilt,omitempty" json:"prebuilt,omitempty" yaml:"prebuilt,omitempty"`

	Params *map[string]any `pkl:"params" bson:"params,omitempty" json:"params,omitempty" yaml:"params,omitempty"`

	WorkingDirectory *string `pkl:"working_directory" bson:"working_directory,omitempty" json:"working_directory,omitempty" yaml:"working_directory,omitempty"`

	Env *map[string]string `pkl:"env" bson:"env,omitempty" json:"env,omitempty" yaml:"env,omitempty"`
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
