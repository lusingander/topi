package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
)

var (
	infoPageTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("70")).
				Bold(true)

	infoPageVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("70"))

	infoPageTitleBarStyle = lipgloss.NewStyle().
				Padding(1, 2)

	infoPageDescriptionStyle = lipgloss.NewStyle().
					Padding(1, 2)

	infoPageUrlStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Underline(true)

	infoPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

type infoPageModel struct {
	doc           *topi.Document
	viewport      viewport.Model
	delegateKeys  infoPageDelegateKeyMap
	width, height int
}

func newInfoPageModel(doc *topi.Document) infoPageModel {
	m := infoPageModel{
		doc: doc,
	}
	m.delegateKeys = newInfoPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	return m
}

type infoPageDelegateKeyMap struct {
	back key.Binding
}

func newInfoPageDelegateKeyMap() infoPageDelegateKeyMap {
	return infoPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *infoPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.viewport.Width, m.viewport.Height = w, h
	m.updateContent()
}

func (m *infoPageModel) updateContent() {
	info := m.doc.Info
	glamourWidth := m.width - 10
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(glamourTheme),
		glamour.WithWordWrap(glamourWidth),
	)

	var content strings.Builder

	title := infoPageTitleStyle.Render(info.Title)
	version := infoPageVersionStyle.Render(fmt.Sprintf("(%s)", info.Version))
	titleBar := infoPageTitleBarStyle.Render(fmt.Sprintf("%s  %s", title, version))
	content.WriteString(titleBar)

	if info.TermsOfService != "" {
		tos := fmt.Sprintf("Terms of service: %s", infoPageUrlStyle.Render(info.TermsOfService))
		content.WriteString(infoPageItemStyle.Render(tos))
	}

	if info.ContactName != "" || info.ContactUrl != "" || info.ContactEmail != "" {
		var buf strings.Builder
		if info.ContactName != "" {
			buf.WriteString(info.ContactName)
			if info.ContactUrl != "" || info.ContactEmail != "" {
				buf.WriteString(": ")
			}
		}
		if info.ContactUrl != "" {
			buf.WriteString(infoPageUrlStyle.Render(info.ContactUrl))
			if info.ContactEmail != "" {
				buf.WriteString(" / ")
			}
		}
		if info.ContactEmail != "" {
			buf.WriteString(infoPageUrlStyle.Render(info.ContactEmail))
		}
		content.WriteString(infoPageItemStyle.Render(buf.String()))
	}

	if info.LicenseName != "" {
		var buf strings.Builder
		buf.WriteString(fmt.Sprintf("License: %s", info.LicenseName))
		if info.LicenseUrl != "" {
			buf.WriteString(fmt.Sprintf(" (%s)", infoPageUrlStyle.Render(info.LicenseUrl)))
		}
		content.WriteString(infoPageItemStyle.Render(buf.String()))
	}

	desc, _ := r.Render(info.Description)
	desc = infoPageDescriptionStyle.Render(desc)
	content.WriteString(desc)

	m.viewport.SetContent(content.String())
}

func (m infoPageModel) Init() tea.Cmd {
	return nil
}

func (m infoPageModel) Update(msg tea.Msg) (infoPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		}
	case selectInfoMenuMsg:
		m.updateContent()
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m infoPageModel) View() string {
	return m.viewport.View()
}
