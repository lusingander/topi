package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type creditsMenuPageListItem struct {
	name        string
	top, bottom int
}

var _ list.Item = (*creditsMenuPageListItem)(nil)

func (i creditsMenuPageListItem) FilterValue() string {
	return ""
}

type creditsMenuPageListDelegate struct{}

var _ list.ItemDelegate = (*creditsMenuPageListDelegate)(nil)

func newCreditsMenuPageListDelegate() creditsMenuPageListDelegate {
	return creditsMenuPageListDelegate{}
}

func (d creditsMenuPageListDelegate) Height() int {
	return 1
}

func (d creditsMenuPageListDelegate) Spacing() int {
	return 0
}

func (d creditsMenuPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d creditsMenuPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(creditsMenuPageListItem)
	selected := index == m.Index()

	name := i.name

	if selected {
		name = listSelectedDescStyle.Render(name)
	} else {
		name = listNormalDescStyle.Render(name)
	}

	fmt.Fprint(w, name)
}
