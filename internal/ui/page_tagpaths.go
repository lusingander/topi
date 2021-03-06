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
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.KeyMap.Quit.Unbind()
	return m
}

type tagPathsPageDelegateKeyMap struct {
	enter key.Binding
	back  key.Binding
}

func newTagPathsPageDelegateKeyMap() tagPathsPageDelegateKeyMap {
	return tagPathsPageDelegateKeyMap{
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

func (m *tagPathsPageModel) reset() {
	m.list.ResetSelected()
	m.list.ResetFilter()
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
		case key.Matches(msg, m.delegateKeys.enter):
			if m.list.FilterState() != list.Filtering {
				path := m.list.SelectedItem().(tagPathsPageListItem).path
				return m, selectOperation(path.OperationId)
			}
		}
	case selectTagMsg:
		m.updateList(msg.tag)
		m.reset()
		return m, nil
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m tagPathsPageModel) View() string {
	return m.list.View()
}

func (m tagPathsPageModel) statusbarInfoString() string {
	return listStatusbarInfoString(m.list)
}

func (m tagPathsPageModel) statusMessageString() string {
	return listStatusMessageString(m.list)
}
