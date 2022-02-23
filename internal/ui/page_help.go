package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	helpPageAboutMenu = "About"
)

var helpPageItems = []list.Item{
	helpPageListItem{
		title:       helpPageAboutMenu,
		description: "Show about this application",
	},
}

type helpPageModel struct {
	list          list.Model
	delegateKeys  helpPageDelegateKeyMap
	width, height int
}

func newHelpPageModel() helpPageModel {
	m := helpPageModel{}
	m.delegateKeys = newHelpPageDelegateKeyMap()
	delegate := newHelpPageListDelegate()
	m.list = list.New(helpPageItems, delegate, 0, 0)
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.SetFilteringEnabled(false) // not necessary on this page...
	m.list.KeyMap.Quit.Unbind()
	return m
}

type helpPageDelegateKeyMap struct {
	back  key.Binding
	enter key.Binding
}

func newHelpPageDelegateKeyMap() helpPageDelegateKeyMap {
	return helpPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *helpPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m helpPageModel) Init() tea.Cmd {
	return nil
}

func (m helpPageModel) Update(msg tea.Msg) (helpPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.enter):
			menu := m.list.SelectedItem().(helpPageListItem)
			switch menu.title {
			case helpPageAboutMenu:
				return m, selectAboutMenu
			}
			return m, nil
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m helpPageModel) View() string {
	return m.list.View()
}
