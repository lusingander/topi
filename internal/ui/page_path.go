package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type pathPageModel struct {
	doc           *topi.Document
	list          list.Model
	delegateKeys  pathPageDelegateKeyMap
	width, height int
}

func newPathPageModel(doc *topi.Document) pathPageModel {
	m := pathPageModel{
		doc: doc,
	}
	m.delegateKeys = newPathPageDelegateKeyMap()
	delegate := newPathPageListDelegate()
	m.list = list.New(nil, delegate, 0, 0)
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	return m
}

type pathPageDelegateKeyMap struct {
	enter key.Binding
	back  key.Binding
}

func newPathPageDelegateKeyMap() pathPageDelegateKeyMap {
	return pathPageDelegateKeyMap{
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *pathPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m *pathPageModel) updateList() {
	m.list.ResetSelected()
	items := make([]list.Item, 0)
	for _, tag := range m.doc.Tags {
		paths := m.doc.TagPathMap[tag.Name]
		for _, path := range paths {
			item := pathPageListItem{path}
			items = append(items, item)
		}
	}
	m.list.SetItems(items)
}

func (m *pathPageModel) reset() {
	m.list.ResetSelected()
	m.list.ResetFilter()
}

func (m pathPageModel) Init() tea.Cmd {
	return nil
}

func (m pathPageModel) Update(msg tea.Msg) (pathPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			if m.list.FilterState() != list.Filtering {
				return m, goBack
			}
		case key.Matches(msg, m.delegateKeys.enter):
			if m.list.FilterState() != list.Filtering {
				path := m.list.SelectedItem().(pathPageListItem).path
				return m, selectOperation(path.OperationId)
			}
		}
	case selectPathMenuMsg:
		m.updateList()
		m.reset()
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m pathPageModel) View() string {
	return m.list.View()
}

func (m pathPageModel) statusbarInfoString() string {
	return listStatusbarInfoString(m.list)
}

func (m pathPageModel) statusMessageString() string {
	return listStatusMessageString(m.list)
}
