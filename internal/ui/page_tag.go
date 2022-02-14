package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type tagPageModel struct {
	list          list.Model
	width, height int
}

func newTagPageModel(doc *topi.Document) tagPageModel {
	m := tagPageModel{}
	items := m.buildItems(doc)
	delegate := newTagPageListDelegate()
	m.list = list.New(items, delegate, 0, 0)
	m.list.Title = topi.AppName
	return m
}

func (tagPageModel) buildItems(doc *topi.Document) []list.Item {
	tags := doc.Tags
	items := make([]list.Item, len(tags))
	for i, tag := range tags {
		item := tagPageListItem{tag}
		items[i] = item
	}
	return items
}

func (m *tagPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m tagPageModel) Init() tea.Cmd {
	return nil
}

func (m tagPageModel) Update(msg tea.Msg) (tagPageModel, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tagPageModel) View() string {
	return m.list.View()
}
