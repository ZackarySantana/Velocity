package uicli

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	choose = key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "choose"),
	)
)

type listModel struct {
	list list.Model
}

func NewListModel(title string, items []list.Item, selected *string) listModel {
	lm := listModel{}
	d := list.NewDefaultDelegate()
	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(SimpleItem); ok {
			title = i.Title()
			*selected = i.Label
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, choose):
				return m.NewStatusMessage(statusMessageStyle("You chose " + title + ". Press q to quit."))
			}
		}

		return nil
	}
	help := []key.Binding{choose}
	d.ShortHelpFunc = func() []key.Binding {
		return help
	}
	d.FullHelpFunc = func() [][]key.Binding {
		return [][]key.Binding{help}
	}
	lm.list = list.New(items, d, 0, 0)
	lm.list.Title = title
	lm.list.Styles.Title = titleStyle
	return lm
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	// This will also call our delegate's update function.
	newListModel, cmd := m.list.Update(msg)
	m.list = newListModel
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m listModel) View() string {
	return appStyle.Render(m.list.View())
}

func Run(title string, items []list.Item) (string, error) {
	k := ""
	selected := &k

	_, err := tea.NewProgram(NewListModel(title, items, selected), tea.WithAltScreen()).Run()
	if err != nil {
		return "", err
	}

	return *selected, nil
}
