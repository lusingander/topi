package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type tagApiPageModel struct {
	doc           *topi.Document
	list          list.Model
	delegateKeys  tagApiPageDelegateKeyMap
	width, height int
}

func newTagApiPageModel(doc *topi.Document) tagApiPageModel {
	m := tagApiPageModel{
		doc: doc,
	}
	m.delegateKeys = newTagApiPageDelegateKeyMap()
	delegate := newTagApiPageListDelegate()
	m.list = list.New(nil, delegate, 0, 0)
	m.list.Title = topi.AppName
	return m
}

type tagApiPageDelegateKeyMap struct {
	back key.Binding
}

func newTagApiPageDelegateKeyMap() tagApiPageDelegateKeyMap {
	return tagApiPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *tagApiPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m *tagApiPageModel) updateList(tag string) {
	m.list.ResetSelected()
	paths := m.doc.TagPathMap[tag]
	items := make([]list.Item, len(paths))
	for i, path := range paths {
		item := tagApiPageListItem{path}
		items[i] = item
	}
	m.list.SetItems(items)
}

func (m tagApiPageModel) Init() tea.Cmd {
	return nil
}

func (m tagApiPageModel) Update(msg tea.Msg) (tagApiPageModel, tea.Cmd) {
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

func (m tagApiPageModel) View() string {
	return m.list.View()
}
