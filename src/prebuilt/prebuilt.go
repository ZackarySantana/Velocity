package prebuilt

import (
	"fmt"

	"github.com/zackarysantana/velocity/src/config/configuration"
)

func GetPrebuilt(prebuilt string) (configuration.PrebuiltCommand, error) {
	switch prebuilt {
	case "git.clone":
		return GitCloneCommand{}, nil
	default:
		return nil, fmt.Errorf("invalid prebuilt command '%s'", prebuilt)
	}
}
