package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var (
	creditsPageContentStyle = lipgloss.NewStyle().
				Padding(0, 2)

	creditsPageRepositoryNameStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("70")).
					Bold(true)

	creditsPageUrlStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Underline(true)

	creditsPageSeparatorColorStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("240"))

	creditsPageSeparator = creditsPageSeparatorColorStyle.
				Render("----------------------------------------")
)

type creditsPageModel struct {
	viewport      viewport.Model
	list          list.Model
	delegateKeys  creditsPageDelegateKeyMap
	width, height int
	showList      bool
}

func newCreditsPageModel() creditsPageModel {
	m := creditsPageModel{}
	m.delegateKeys = newCreditsPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	items := make([]list.Item, len(credits))
	for i, c := range credits {
		item := creditsMenuPageListItem{name: c.name}
		items[i] = item
	}
	m.list = list.New(items, newCreditsMenuPageListDelegate(), 0, 0)
	m.list.SetShowTitle(false)
	m.list.SetShowFilter(false)
	m.list.SetShowPagination(false)
	m.list.SetShowStatusBar(false)
	m.list.SetShowHelp(false)
	m.list.SetFilteringEnabled(false)
	return m
}

type creditsPageDelegateKeyMap struct {
	back   key.Binding
	toggle key.Binding
}

func newCreditsPageDelegateKeyMap() creditsPageDelegateKeyMap {
	return creditsPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
		toggle: key.NewBinding(
			key.WithKeys("t"),
			key.WithHelp("t", "toggle"),
		),
	}
}

func (m *creditsPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	if m.showList {
		if h%2 == 0 {
			m.list.SetSize(w, h/2-1)
		} else {
			m.list.SetSize(w, h/2)
		}
		m.viewport.Width, m.viewport.Height = w, h/2
	} else {
		m.list.SetSize(0, 0)
		m.viewport.Width, m.viewport.Height = w, h
	}
	m.updateContent()
}

func (m *creditsPageModel) reset() {
	m.viewport.GotoTop()
	m.list.ResetSelected()
	m.showList = false
	m.SetSize(m.width, m.height)
}

func (m *creditsPageModel) updateContent() {
	var content strings.Builder
	n := 0
	for i, credit := range credits {
		name := creditsPageRepositoryNameStyle.Render(credit.name)
		url := creditsPageUrlStyle.Render(credit.url)
		text := wordwrap.String(credit.text, m.width)

		s := fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n", name, url, text, creditsPageSeparator)
		content.WriteString(s)

		item := m.list.Items()[i].(creditsMenuPageListItem)
		item.top = n
		n += strings.Count(s, "\n")
		item.bottom = n - 1
		m.list.SetItem(i, item)
	}
	m.viewport.SetContent(creditsPageContentStyle.Render(content.String()))
}

func (m *creditsPageModel) toggleList() {
	m.showList = !m.showList
	m.SetSize(m.width, m.height)
}

func (m *creditsPageModel) updateViewportPos() {
	item := m.list.SelectedItem().(creditsMenuPageListItem)
	m.viewport.SetYOffset(item.top)
}

func (m *creditsPageModel) updateListPos() {
	pos := m.viewport.YOffset
	for i, item := range m.list.Items() {
		cItem := item.(creditsMenuPageListItem)
		if cItem.top <= pos && pos < cItem.bottom {
			m.list.Select(i)
			return
		}
	}
}

func (m creditsPageModel) Init() tea.Cmd {
	return nil
}

func (m creditsPageModel) Update(msg tea.Msg) (creditsPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		case key.Matches(msg, m.delegateKeys.toggle):
			m.toggleList()
			return m, nil
		}
	case selectCreditsMenuMsg:
		m.reset()
		m.updateContent()
		return m, nil
	}

	var cmd tea.Cmd
	if m.showList {
		m.list, cmd = m.list.Update(msg)
		m.updateViewportPos()
	} else {
		m.viewport, cmd = m.viewport.Update(msg)
		m.updateListPos()
	}
	return m, cmd
}

func (m creditsPageModel) View() string {
	v := m.viewport.View()
	if !m.showList {
		return v
	}
	l := m.list.View()
	sep := creditsPageSeparatorColorStyle.Render(strings.Repeat("-", m.width))
	return lipgloss.JoinVertical(lipgloss.Top, v, sep, l)
}
