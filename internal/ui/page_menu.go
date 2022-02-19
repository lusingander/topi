package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

const (
	menuPageInfoMenu  = "Info"
	menuPagePathsMenu = "Paths"
	menuPageHelpMenu  = "Help"
)

var menuPageItems = []list.Item{
	menuPageListItem{
		title:       menuPageInfoMenu,
		description: "Show API information",
	},
	menuPageListItem{
		title:       menuPagePathsMenu,
		description: "Show all paths (grouped by tags)",
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
	m.list.Title = topi.AppName
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
			case menuPagePathsMenu:
				return m, selectTagMenu
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
