package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type helpMenuPageListItem struct {
	title       string
	description string
}

var _ list.Item = (*helpMenuPageListItem)(nil)

func (i helpMenuPageListItem) FilterValue() string {
	return ""
}

type helpMenuPageListDelegate struct{}

var _ list.ItemDelegate = (*helpMenuPageListDelegate)(nil)

func newHelpMenuPageListDelegate() helpMenuPageListDelegate {
	return helpMenuPageListDelegate{}
}

func (d helpMenuPageListDelegate) Height() int {
	return 2
}

func (d helpMenuPageListDelegate) Spacing() int {
	return 1
}

func (d helpMenuPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d helpMenuPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(helpMenuPageListItem)
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
