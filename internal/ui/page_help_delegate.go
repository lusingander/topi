package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type helpPageListItem struct {
	title       string
	description string
}

var _ list.Item = (*helpPageListItem)(nil)

func (i helpPageListItem) FilterValue() string {
	return ""
}

type helpPageListDelegate struct{}

var _ list.ItemDelegate = (*helpPageListDelegate)(nil)

func newHelpPageListDelegate() helpPageListDelegate {
	return helpPageListDelegate{}
}

func (d helpPageListDelegate) Height() int {
	return 2
}

func (d helpPageListDelegate) Spacing() int {
	return 1
}

func (d helpPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d helpPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(helpPageListItem)
	selected := index == m.Index()

	title := i.title
	desc := i.description

	if selected {
		title = listSelectedTitleStyle.Render(title)
		desc = listSelectedDescStyle.Render(desc)
	} else {
		title = listNormalTitleStyle.Render(title)
		desc = listNormalDescStyle.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
