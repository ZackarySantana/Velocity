// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

type DockerRuntime interface {
	Runtime

	GetImage() string
}

var _ DockerRuntime = (*DockerRuntimeImpl)(nil)

type DockerRuntimeImpl struct {
	Image string `pkl:"image" bson:"image,omitempty" json:"image,omitempty" yaml:"image,omitempty"`

	Name string `pkl:"name" bson:"name,omitempty" json:"name,omitempty" yaml:"name,omitempty"`

	Env *map[string]string `pkl:"env" bson:"env,omitempty" json:"env,omitempty" yaml:"env,omitempty"`
}

func (rcv *DockerRuntimeImpl) GetImage() string {
	return rcv.Image
}

func (rcv *DockerRuntimeImpl) GetName() string {
	return rcv.Name
}

func (rcv *DockerRuntimeImpl) GetEnv() *map[string]string {
	return rcv.Env
}
