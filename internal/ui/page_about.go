package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
)

var (
	aboutPageTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("70")).
				Bold(true)

	aboutPageTitleBarStyle = lipgloss.NewStyle().
				Padding(0, 2, 1)

	aboutPageUrlStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Underline(true)

	aboutPageSelectedUrlStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("250")).
					Foreground(lipgloss.Color("56"))

	aboutPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

type aboutPageSelectableItems int

const (
	aboutPageSelectableNotSelected aboutPageSelectableItems = iota
	aboutPageSelectableAppUrl
	aboutPageSelectableNumberOfItems // not item
)

type aboutPageModel struct {
	viewport      viewport.Model
	delegateKeys  aboutPageDelegateKeyMap
	width, height int

	selected aboutPageSelectableItems
}

func newAboutPageModel() aboutPageModel {
	m := aboutPageModel{}
	m.delegateKeys = newAboutPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	return m
}

type aboutPageDelegateKeyMap struct {
	back        key.Binding
	tab         key.Binding
	shiftTab    key.Binding
	openBrowser key.Binding
}

func newAboutPageDelegateKeyMap() aboutPageDelegateKeyMap {
	return aboutPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
		tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "select next item"),
		),
		shiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "select prev item"),
		),
		openBrowser: key.NewBinding(
			key.WithKeys("x"),
			key.WithHelp("x", "open in browser"),
		),
	}
}

func (m *aboutPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.viewport.Width, m.viewport.Height = w, h
	m.updateContent()
}

func (m *aboutPageModel) reset() {
	m.selected = aboutPageSelectableNotSelected
	m.viewport.GotoTop()
}

func (m *aboutPageModel) updateContent() {
	var content strings.Builder

	title := aboutPageTitleStyle.Render(topi.AppName)
	titleBar := aboutPageTitleBarStyle.Render(title)
	content.WriteString(titleBar)

	version := aboutPageItemStyle.Render(fmt.Sprintf("Version %s", topi.AppVersion))
	content.WriteString(version)

	url := topi.AppUrl
	if m.selected == aboutPageSelectableAppUrl {
		url = aboutPageSelectedUrlStyle.Render(url)
	} else {
		url = aboutPageUrlStyle.Render(url)
	}
	content.WriteString(aboutPageItemStyle.Render(url))

	m.viewport.SetContent(content.String())
}

func (m *aboutPageModel) selectItem(reverse bool) {
	n := aboutPageSelectableNumberOfItems
	if reverse {
		m.selected = ((m.selected-1)%n + n) % n
	} else {
		m.selected = (m.selected + 1) % n
	}
	switch m.selected {
	case aboutPageSelectableAppUrl:
		// do nothing
	default:
		m.selected = aboutPageSelectableNotSelected
	}
}

func (m aboutPageModel) openInBrowser() error {
	switch m.selected {
	case aboutPageSelectableAppUrl:
		return openInBrowser(topi.AppUrl)
	default:
		return nil // do nothing
	}
}

func (m aboutPageModel) Init() tea.Cmd {
	return nil
}

func (m aboutPageModel) Update(msg tea.Msg) (aboutPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		case key.Matches(msg, m.delegateKeys.tab):
			m.selectItem(false)
			m.updateContent()
			return m, nil
		case key.Matches(msg, m.delegateKeys.shiftTab):
			m.selectItem(true)
			m.updateContent()
			return m, nil
		case key.Matches(msg, m.delegateKeys.openBrowser):
			m.openInBrowser() // todo: handle error
			return m, nil
		}
	case selectAboutMenuMsg:
		m.reset()
		m.updateContent()
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m aboutPageModel) View() string {
	return m.viewport.View()
}
