package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	helpMenuPageHelpMenu  = "Help"
	helpMenuPageAboutMenu = "About"
)

var helpMenuPageItems = []list.Item{
	helpMenuPageListItem{
		title:       helpMenuPageHelpMenu,
		description: "Show help",
	},
	helpMenuPageListItem{
		title:       helpMenuPageAboutMenu,
		description: "Show about this application",
	},
}

type helpMenuPageModel struct {
	list          list.Model
	delegateKeys  helpMenuPageDelegateKeyMap
	width, height int
}

func newHelpMenuPageModel() helpMenuPageModel {
	m := helpMenuPageModel{}
	m.delegateKeys = newHelpMenuPageDelegateKeyMap()
	delegate := newHelpMenuPageListDelegate()
	m.list = list.New(helpMenuPageItems, delegate, 0, 0)
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.SetFilteringEnabled(false) // not necessary on this page...
	m.list.KeyMap.Quit.Unbind()
	return m
}

type helpMenuPageDelegateKeyMap struct {
	back  key.Binding
	enter key.Binding
}

func newHelpMenuPageDelegateKeyMap() helpMenuPageDelegateKeyMap {
	return helpMenuPageDelegateKeyMap{
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

func (m *helpMenuPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m helpMenuPageModel) Init() tea.Cmd {
	return nil
}

func (m helpMenuPageModel) Update(msg tea.Msg) (helpMenuPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.enter):
			menu := m.list.SelectedItem().(helpMenuPageListItem)
			switch menu.title {
			case helpMenuPageHelpMenu:
				return m, selectHelpHelpMenu
			case helpMenuPageAboutMenu:
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

func (m helpMenuPageModel) View() string {
	return m.list.View()
}
