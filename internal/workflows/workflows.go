package workflows

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/uicli"
)

func GetWorkflow(c config.Config, title string) (config.YAMLWorkflow, error) {
	items := []list.Item{}
	for title, workflow := range c.Workflows {
		desc := workflow.Description
		if desc == nil {
			descLiteral := ""
			desc = &descLiteral
		}
		items = append(items, uicli.SimpleItem{
			Label: title,
			Desc:  *desc,
		})
	}

	result, err := uicli.Run(title, items)
	if err != nil {
		return config.YAMLWorkflow{}, err
	}
	for title, workflow := range c.Workflows {
		if title == result {
			return workflow, nil
		}
	}

	return config.YAMLWorkflow{}, fmt.Errorf("workflow selected not found")
}
