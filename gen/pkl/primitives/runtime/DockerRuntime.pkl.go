// Code generated from Pkl module `runtime_primitive`. DO NOT EDIT.
package runtime

type DockerRuntime interface {
	Runtime

	GetImage() string
}

var _ DockerRuntime = (*DockerRuntimeImpl)(nil)

type DockerRuntimeImpl struct {
	Image string `pkl:"image"`

	Name string `pkl:"name"`

	Env *map[string]string `pkl:"env"`
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
