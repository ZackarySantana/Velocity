package prebuilt

import (
	"github.com/zackarysantana/velocity/src/config/configuration"
)

type GitCloneCommand struct {
	Info         configuration.CommandInfo
	PrebuiltInfo PrebuiltInfo

	Prebuilt_ string
	Params_   []map[string]string
}

func NewGitClone(info configuration.CommandInfo, prebuiltInfo PrebuiltInfo) configuration.PrebuiltCommand {
	return GitCloneCommand{
		Info:         info,
		PrebuiltInfo: prebuiltInfo,
		Prebuilt_:    "git-clone",
	}
}

func (g GitCloneCommand) GetInfo() configuration.CommandInfo {
	return g.Info
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
