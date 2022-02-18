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

	infoPageSelectedUrlStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("250")).
					Foreground(lipgloss.Color("56"))

	infoPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

type infoPageSelectableItems int

const (
	infoPageSelectableNotSelected infoPageSelectableItems = iota
	infoPageSelectableTermsOfService
	infoPageSelectableContractUrl
	infoPageSelectableLicenseUrl
	numberOfItems // not item
)

type infoPageModel struct {
	doc           *topi.Document
	viewport      viewport.Model
	delegateKeys  infoPageDelegateKeyMap
	width, height int

	selected infoPageSelectableItems
}

func newInfoPageModel(doc *topi.Document) infoPageModel {
	m := infoPageModel{
		doc: doc,
	}
	m.delegateKeys = newInfoPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	m.selected = infoPageSelectableNotSelected
	return m
}

type infoPageDelegateKeyMap struct {
	back        key.Binding
	tab         key.Binding
	shiftTab    key.Binding
	openBrowser key.Binding
}

func newInfoPageDelegateKeyMap() infoPageDelegateKeyMap {
	return infoPageDelegateKeyMap{
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
		url := info.TermsOfService
		if m.selected == infoPageSelectableTermsOfService {
			url = infoPageSelectedUrlStyle.Render(url)
		} else {
			url = infoPageUrlStyle.Render(url)
		}
		tos := fmt.Sprintf("Terms of service: %s", url)
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
			url := info.ContactUrl
			if m.selected == infoPageSelectableContractUrl {
				url = infoPageSelectedUrlStyle.Render(url)
			} else {
				url = infoPageUrlStyle.Render(url)
			}
			buf.WriteString(url)
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
			url := info.LicenseUrl
			if m.selected == infoPageSelectableLicenseUrl {
				url = infoPageSelectedUrlStyle.Render(url)
			} else {
				url = infoPageUrlStyle.Render(url)
			}
			buf.WriteString(fmt.Sprintf(" (%s)", url))
		}
		content.WriteString(infoPageItemStyle.Render(buf.String()))
	}

	desc, _ := r.Render(info.Description)
	desc = infoPageDescriptionStyle.Render(desc)
	content.WriteString(desc)

	m.viewport.SetContent(content.String())
}

func (m *infoPageModel) selectItem(reverse bool) {
	if reverse {
		m.selected = ((m.selected-1)%numberOfItems + numberOfItems) % numberOfItems
	} else {
		m.selected = (m.selected + 1) % numberOfItems
	}
	switch m.selected {
	case infoPageSelectableTermsOfService:
		if m.doc.Info.TermsOfService == "" {
			m.selectItem(reverse)
		}
	case infoPageSelectableContractUrl:
		if m.doc.Info.ContactUrl == "" {
			m.selectItem(reverse)
		}
	case infoPageSelectableLicenseUrl:
		if m.doc.Info.LicenseUrl == "" {
			m.selectItem(reverse)
		}
	default:
		m.selected = infoPageSelectableNotSelected
	}
}

func (m infoPageModel) openInBrowser() error {
	switch m.selected {
	case infoPageSelectableTermsOfService:
		return openInBrowser(m.doc.Info.TermsOfService)
	case infoPageSelectableContractUrl:
		return openInBrowser(m.doc.Info.ContactUrl)
	case infoPageSelectableLicenseUrl:
		return openInBrowser(m.doc.Info.LicenseUrl)
	default:
		return nil // do nothing
	}
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
