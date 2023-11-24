package uicli

type SimpleItem struct {
	Label string
	Desc  string
}

func (i SimpleItem) Title() string       { return i.Label }
func (i SimpleItem) Description() string { return i.Desc }
func (i SimpleItem) FilterValue() string { return i.Label }
