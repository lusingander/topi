package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type tagPageListItem struct {
	tag *topi.Tag
}

var _ list.Item = (*tagPageListItem)(nil)

func (i tagPageListItem) FilterValue() string {
	return i.tag.Name
}

func (i tagPageListItem) desc() string {
	if i.tag.Description == "" {
		return "-"
	}
	return i.tag.Description
}

type tagPageListDelegate struct{}

var _ list.ItemDelegate = (*tagPageListDelegate)(nil)

func newTagPageListDelegate() tagPageListDelegate {
	return tagPageListDelegate{}
}

func (d tagPageListDelegate) Height() int {
	return 2
}

func (d tagPageListDelegate) Spacing() int {
	return 1
}

func (d tagPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d tagPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(tagPageListItem)
	selected := index == m.Index()

	title := i.tag.Name
	desc := i.desc()

	if selected {
		title = listSelectedTitleStyle.Render(title)
		desc = listSelectedDescStyle.Render(desc)
	} else {
		title = listNormalTitleStyle.Render(title)
		desc = listNormalDescStyle.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
