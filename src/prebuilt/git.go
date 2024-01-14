package prebuilt

import (
	"github.com/zackarysantana/velocity/src/config/configuration"
	"github.com/zackarysantana/velocity/src/env"
)

type GitCloneCommand struct {
	WorkingDirectory_ *string
	Env_              *env.Env

	Prebuilt_ string
	Params_   []map[string]string
}

func NewGitClone(wd *string, env *env.Env, params []map[string]string) GitCloneCommand {
	return GitCloneCommand{
		WorkingDirectory_: wd,
		Env_:              env,
		Prebuilt_:         "git-clone",
		Params_:           params,
	}
}

func (g GitCloneCommand) WorkingDirectory() *string {
	return g.WorkingDirectory_
}

func (g GitCloneCommand) Env() *env.Env {
	return g.Env_
}

func (g GitCloneCommand) Prebuilt() string {
	return g.Prebuilt_
}

func (g GitCloneCommand) Params() []map[string]string {
	return g.Params_
}

func (g GitCloneCommand) Validate(c configuration.Configuration) error {
	return nil
}
