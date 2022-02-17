package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type menuPageListItem struct {
	title       string
	description string
}

var _ list.Item = (*menuPageListItem)(nil)

func (i menuPageListItem) FilterValue() string {
	return ""
}

type menuPageListDelegate struct{}

var _ list.ItemDelegate = (*menuPageListDelegate)(nil)

func newMenuPageListDelegate() menuPageListDelegate {
	return menuPageListDelegate{}
}

func (d menuPageListDelegate) Height() int {
	return 2
}

func (d menuPageListDelegate) Spacing() int {
	return 1
}

func (d menuPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d menuPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(menuPageListItem)
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
