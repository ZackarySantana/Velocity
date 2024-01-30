package prebuilt

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config/configuration"
)

type PrebuiltInfo struct {
	Params []map[string]string
}

func GetPrebuilt(prebuilt string) (configuration.PrebuiltCommand, error) {
	switch prebuilt {
	case "git.clone":
		return GitCloneCommand{}, nil
	default:
		return nil, fmt.Errorf("invalid prebuilt command '%s'", prebuilt)
	}
}

func GetPrebuiltConstructor(prebuilt string) (func(configuration.CommandInfo, PrebuiltInfo) configuration.PrebuiltCommand, error) {
	switch prebuilt {
	case "git.clone":
		return NewGitClone, nil
	default:
		return nil, fmt.Errorf("invalid prebuilt command '%s'", prebuilt)
	}
}
