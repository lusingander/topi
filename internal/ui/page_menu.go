package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	menuPageInfoMenu  = "Info"
	menuPageTagsMenu  = "Tags"
	menuPagePathsMenu = "Paths"
	menuPageHelpMenu  = "Help"
)

var menuPageItems = []list.Item{
	menuPageListItem{
		title:       menuPageInfoMenu,
		description: "Show API information",
	},
	menuPageListItem{
		title:       menuPageTagsMenu,
		description: "Show all tags",
	},
	menuPageListItem{
		title:       menuPagePathsMenu,
		description: "Show all paths",
	},
	menuPageListItem{
		title:       menuPageHelpMenu,
		description: "Show help menus",
	},
}

type menuPageModel struct {
	list          list.Model
	delegateKeys  menuPageDelegateKeyMap
	width, height int
}

func newMenuPageModel() menuPageModel {
	m := menuPageModel{}
	m.delegateKeys = newMenuPageDelegateKeyMap()
	delegate := newMenuPageListDelegate()
	m.list = list.New(menuPageItems, delegate, 0, 0)
	m.list.SetShowTitle(false)
	m.list.SetShowHelp(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.SetFilteringEnabled(false) // not necessary on this page...
	m.list.KeyMap.Quit.Unbind()
	return m
}

type menuPageDelegateKeyMap struct {
	enter key.Binding
}

func newMenuPageDelegateKeyMap() menuPageDelegateKeyMap {
	return menuPageDelegateKeyMap{
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

func (m *menuPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.list.SetSize(w, h)
}

func (m menuPageModel) Init() tea.Cmd {
	return nil
}

func (m menuPageModel) Update(msg tea.Msg) (menuPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.enter):
			menu := m.list.SelectedItem().(menuPageListItem)
			switch menu.title {
			case menuPageInfoMenu:
				return m, selectInfoMenu
			case menuPageTagsMenu:
				return m, selectTagMenu
			case menuPagePathsMenu:
				return m, selectPathMenu
			case menuPageHelpMenu:
				return m, selectHelpMenu
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuPageModel) View() string {
	return m.list.View()
}
