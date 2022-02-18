package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type tagPageModel struct {
	doc           *topi.Document
	list          list.Model
	delegateKeys  tagPageDelegateKeyMap
	width, height int
}

func newTagPageModel(doc *topi.Document) tagPageModel {
	m := tagPageModel{
		doc: doc,
	}
	m.delegateKeys = newTagPageDelegateKeyMap()
	delegate := newTagPageListDelegate()
	m.list = list.New(nil, delegate, 0, 0)
	m.list.Title = topi.AppName
	return m
}

type tagPageDelegateKeyMap struct {
	back key.Binding
	sel  key.Binding
}

func newTagPageDelegateKeyMap() tagPageDelegateKeyMap {
	return tagPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
		sel: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *tagPageModel) updateItems() {
	tags := m.doc.Tags
	items := make([]list.Item, 0)
	for _, tag := range tags {
		if len(m.doc.TagPathMap[tag.Name]) > 0 {
			item := tagPageListItem{tag}
			items = append(items, item)
		}
	}
	m.list.SetItems(items)
}

func (m *tagPageModel) reset() {
	m.list.ResetSelected()
}

func (m *tagPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m tagPageModel) Init() tea.Cmd {
	return nil
}

func (m tagPageModel) Update(msg tea.Msg) (tagPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		case key.Matches(msg, m.delegateKeys.sel):
			tag := m.list.SelectedItem().(tagPageListItem).tag
			return m, selectTag(tag.Name)
		}
	case selectTagMenuMsg:
		m.updateItems()
		m.reset()
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tagPageModel) View() string {
	return m.list.View()
}
