package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

	creditsPageSeparator = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("----------------------------------------")
)

type creditsPageModel struct {
	viewport      viewport.Model
	delegateKeys  creditsPageDelegateKeyMap
	width, height int
}

func newCreditsPageModel() creditsPageModel {
	m := creditsPageModel{}
	m.delegateKeys = newCreditsPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	return m
}

type creditsPageDelegateKeyMap struct {
	back key.Binding
}

func newCreditsPageDelegateKeyMap() creditsPageDelegateKeyMap {
	return creditsPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *creditsPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.viewport.Width, m.viewport.Height = w, h
	m.updateContent()
}

func (m *creditsPageModel) reset() {
	m.viewport.GotoTop()
}

func (m *creditsPageModel) updateContent() {
	var content strings.Builder
	for _, credit := range credits {
		name := creditsPageRepositoryNameStyle.Render(credit.name)
		url := creditsPageUrlStyle.Render(credit.url)
		text := credit.text
		content.WriteString(fmt.Sprintf("%s\n\n%s\n\n%s\n%s\n\n", name, url, text, creditsPageSeparator))
	}
	m.viewport.SetContent(creditsPageContentStyle.Render(content.String()))
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
		}
	case selectCreditsMenuMsg:
		m.reset()
		m.updateContent()
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m creditsPageModel) View() string {
	return m.viewport.View()
}
