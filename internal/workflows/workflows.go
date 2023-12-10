package workflows

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/zackarysantana/velocity/src/config"
	"github.com/zackarysantana/velocity/src/uicli"
)

func GetWorkflow(c *config.Config, title string) (config.YAMLWorkflow, error) {
	items := []uicli.SimpleItem{}
	for name, workflow := range c.Workflows {
		desc := workflow.Description
		if desc == nil {
			descLiteral := ""
			desc = &descLiteral
		}
		items = append(items, uicli.SimpleItem{
			Label: name,
			Desc:  *desc,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return strings.Compare(items[i].Label, items[j].Label) < 0
	})

	var listItems []list.Item
	for _, item := range items {
		listItems = append(listItems, item)
	}

	result, err := uicli.Run(title, listItems)
	if err != nil {
		return config.YAMLWorkflow{}, err
	}
	for name, workflow := range c.Workflows {
		if name == result {
			return workflow, nil
		}
	}

	return config.YAMLWorkflow{}, fmt.Errorf("workflow selected not found")
}
