package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type tagPathsPageModel struct {
	doc           *topi.Document
	list          list.Model
	delegateKeys  tagPathsPageDelegateKeyMap
	width, height int
}

func newTagPathsPageModel(doc *topi.Document) tagPathsPageModel {
	m := tagPathsPageModel{
		doc: doc,
	}
	m.delegateKeys = newTagPathsPageDelegateKeyMap()
	delegate := newTagPathsPageListDelegate()
	m.list = list.New(nil, delegate, 0, 0)
	m.list.Title = topi.AppName
	return m
}

type tagPathsPageDelegateKeyMap struct {
	back key.Binding
}

func newTagPathsPageDelegateKeyMap() tagPathsPageDelegateKeyMap {
	return tagPathsPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *tagPathsPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m *tagPathsPageModel) updateList(tag string) {
	m.list.ResetSelected()
	paths := m.doc.TagPathMap[tag]
	items := make([]list.Item, len(paths))
	for i, path := range paths {
		item := tagPathsPageListItem{path}
		items[i] = item
	}
	m.list.SetItems(items)
}

func (m tagPathsPageModel) Init() tea.Cmd {
	return nil
}

func (m tagPathsPageModel) Update(msg tea.Msg) (tagPathsPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			if m.list.FilterState() != list.Filtering {
				return m, goBack
			}
		}
	case selectTagMsg:
		m.updateList(msg.tag)
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tagPathsPageModel) View() string {
	return m.list.View()
}
