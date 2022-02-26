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
	infoPageTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("70")).
				Bold(true)

	infoPageVersionStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("70"))

	infoPageTitleBarStyle = lipgloss.NewStyle().
				Padding(0, 2, 1)

	infoPageDescriptionStyle = lipgloss.NewStyle().
					Padding(1, 2)

	infoPageUrlStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("33")).
				Underline(true)

	infoPageSelectedUrlStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("250")).
					Foreground(lipgloss.Color("56"))

	infoPageSectionHeaderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("70")).
					Underline(true)

	infoPageSectionSubHeaderStyle = infoPageSectionHeaderStyle.Copy().
					Margin(0, 0, 0, 1)

	infoPageAuthenticationItemStyle = infoPageItemStyle.Copy().
					Margin(0, 0, 0, 2)

	infoPageAuthenticationItemDescriptionStyle = infoPageAuthenticationItemStyle.Copy().
							PaddingBottom(0)

	infoPageAuthenticationItemKeyColorStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("143"))

	infoPageAuthenticationItemValueColorStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("167"))

	infoPageAuthenticationOAuthFlowStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("70"))

	infoPageAuthenticationOAuthScopesStyle = lipgloss.NewStyle().
						MarginLeft(2)

	infoPageAuthenticationOAuthScopeNameColorStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("167")).
							Underline(true)

	infoPageAuthenticationOAuthScopeDescColorStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("246"))

	infoPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

var (
	infoPageSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(1, 2).
		Render("----------")
)

type infoPageSelectableItems int

const (
	infoPageSelectableNotSelected infoPageSelectableItems = iota
	infoPageSelectableTermsOfService
	infoPageSelectableContractUrl
	infoPageSelectableLicenseUrl
	infoPageSelectableExDocsUrl
	infoPageSelectableNumberOfItems // not item
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

func (m *infoPageModel) reset() {
	m.selected = infoPageSelectableNotSelected
	m.viewport.GotoTop()
}

func (m *infoPageModel) updateContent() {
	info := m.doc.Info
	r, _ := markdownRenderer(m.width - 10)

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

	if info.Description != "" {
		desc, _ := r.Render(info.Description)
		desc = infoPageDescriptionStyle.Render(desc)
		content.WriteString(desc)
	}

	if info.ExDocsDescription != "" {
		desc, _ := r.Render(info.ExDocsDescription)
		desc = infoPageDescriptionStyle.Render(desc)
		content.WriteString(desc)
	}

	if info.ExDocsUrl != "" {
		url := info.ExDocsUrl
		if m.selected == infoPageSelectableExDocsUrl {
			url = infoPageSelectedUrlStyle.Render(url)
		} else {
			url = infoPageUrlStyle.Render(url)
		}
		content.WriteString(infoPageItemStyle.Render(url))
	}

	if m.doc.Components != nil {
		schemes := m.doc.Components.SecuritySchemes
		if len(schemes) > 0 {
			h := infoPageSectionHeaderStyle.Render("Authentication")
			content.WriteString(infoPageItemStyle.Render(h))

			for _, scheme := range schemes {
				h := infoPageSectionSubHeaderStyle.Render(fmt.Sprintf("%s (%s)", scheme.Key, scheme.TypeStr()))
				content.WriteString(infoPageItemStyle.Render(h))

				if scheme.Description != "" {
					desc, _ := r.Render(scheme.Description)
					desc = infoPageAuthenticationItemDescriptionStyle.Render(desc)
					content.WriteString(desc)
				}

				switch scheme.Type {
				case "apiKey":
					nameKey := infoPageAuthenticationItemKeyColorStyle.Render("Parameter name:")
					nameValue := infoPageAuthenticationItemValueColorStyle.Render(scheme.Name)
					inKey := infoPageAuthenticationItemKeyColorStyle.Render("Parameter in:")
					inValue := infoPageAuthenticationItemValueColorStyle.Render(scheme.In)
					values := fmt.Sprintf("%s %s\n%s %s", nameKey, nameValue, inKey, inValue)
					content.WriteString(infoPageAuthenticationItemStyle.Render(values))
				case "http":
					schemeKey := infoPageAuthenticationItemKeyColorStyle.Render("HTTP Authorization Scheme:")
					schemeValue := infoPageAuthenticationItemValueColorStyle.Render(scheme.Scheme)
					values := fmt.Sprintf("%s %s", schemeKey, schemeValue)
					if scheme.BearerFormat != "" {
						formatKey := infoPageAuthenticationItemKeyColorStyle.Render("Bearer format:")
						formatValue := infoPageAuthenticationItemValueColorStyle.Render(scheme.BearerFormat)
						values += fmt.Sprintf("\n%s %s", formatKey, formatValue)
					}
					content.WriteString(infoPageAuthenticationItemStyle.Render(values))
				case "oauth2":
					if scheme.OAuthFlows.AuthorizatonCodeFlow != nil {
						f := scheme.OAuthFlows.AuthorizatonCodeFlow
						flow := infoPageAuthenticationOAuthFlowStyle.Render("[Authorization Code Flow]")
						content.WriteString(infoPageAuthenticationItemStyle.Render(flow))
						authKey := infoPageAuthenticationItemKeyColorStyle.Render("Authorization URL:")
						authValue := infoPageAuthenticationItemValueColorStyle.Render(f.AuthorizationURL)
						tokenKey := infoPageAuthenticationItemKeyColorStyle.Render("Token URL:")
						tokenValue := infoPageAuthenticationItemValueColorStyle.Render(f.TokenURL)
						values := fmt.Sprintf("%s %s\n%s %s", authKey, authValue, tokenKey, tokenValue)
						if f.RefreshURL != "" {
							refKey := infoPageAuthenticationItemKeyColorStyle.Render("Reflesh URL:")
							refValue := infoPageAuthenticationItemValueColorStyle.Render(f.RefreshURL)
							values += fmt.Sprintf("\n%s %s", refKey, refValue)
						}
						scopesKey := infoPageAuthenticationItemKeyColorStyle.Render("Scopes:")
						scopesValue := infoPageAuthenticationOAuthScopesStyle.Render(styledOAuthScopes(f.Scopes))
						values += fmt.Sprintf("\n%s\n%s", scopesKey, scopesValue)
						content.WriteString(infoPageAuthenticationItemStyle.Render(values))
					}
					if scheme.OAuthFlows.ImplicitFlow != nil {
						f := scheme.OAuthFlows.ImplicitFlow
						flow := infoPageAuthenticationOAuthFlowStyle.Render("[Implicit Flow]")
						content.WriteString(infoPageAuthenticationItemStyle.Render(flow))
						authKey := infoPageAuthenticationItemKeyColorStyle.Render("Authorization URL:")
						authValue := infoPageAuthenticationItemValueColorStyle.Render(f.AuthorizationURL)
						values := fmt.Sprintf("%s %s", authKey, authValue)
						if f.RefreshURL != "" {
							refKey := infoPageAuthenticationItemKeyColorStyle.Render("Reflesh URL:")
							refValue := infoPageAuthenticationItemValueColorStyle.Render(f.RefreshURL)
							values += fmt.Sprintf("\n%s %s", refKey, refValue)
						}
						scopesKey := infoPageAuthenticationItemKeyColorStyle.Render("Scopes:")
						scopesValue := infoPageAuthenticationOAuthScopesStyle.Render(styledOAuthScopes(f.Scopes))
						values += fmt.Sprintf("\n%s\n%s", scopesKey, scopesValue)
						content.WriteString(infoPageAuthenticationItemStyle.Render(values))
					}
					if scheme.OAuthFlows.ResourceOwnerPasswordCredentialsFlow != nil {
						f := scheme.OAuthFlows.ResourceOwnerPasswordCredentialsFlow
						flow := infoPageAuthenticationOAuthFlowStyle.Render("[Resource Owner Password Credentials Flow]")
						content.WriteString(infoPageAuthenticationItemStyle.Render(flow))
						tokenKey := infoPageAuthenticationItemKeyColorStyle.Render("Token URL:")
						tokenValue := infoPageAuthenticationItemValueColorStyle.Render(f.TokenURL)
						values := fmt.Sprintf("%s %s", tokenKey, tokenValue)
						if f.RefreshURL != "" {
							refKey := infoPageAuthenticationItemKeyColorStyle.Render("Reflesh URL:")
							refValue := infoPageAuthenticationItemValueColorStyle.Render(f.RefreshURL)
							values += fmt.Sprintf("\n%s %s", refKey, refValue)
						}
						scopesKey := infoPageAuthenticationItemKeyColorStyle.Render("Scopes:")
						scopesValue := infoPageAuthenticationOAuthScopesStyle.Render(styledOAuthScopes(f.Scopes))
						values += fmt.Sprintf("\n%s\n%s", scopesKey, scopesValue)
						content.WriteString(infoPageAuthenticationItemStyle.Render(values))
					}
					if scheme.OAuthFlows.ClientCredentialsFlow != nil {
						f := scheme.OAuthFlows.ClientCredentialsFlow
						flow := infoPageAuthenticationOAuthFlowStyle.Render("[Client Credentials Flow]")
						content.WriteString(infoPageAuthenticationItemStyle.Render(flow))
						tokenKey := infoPageAuthenticationItemKeyColorStyle.Render("Token URL:")
						tokenValue := infoPageAuthenticationItemValueColorStyle.Render(f.TokenURL)
						values := fmt.Sprintf("%s %s", tokenKey, tokenValue)
						if f.RefreshURL != "" {
							refKey := infoPageAuthenticationItemKeyColorStyle.Render("Reflesh URL:")
							refValue := infoPageAuthenticationItemValueColorStyle.Render(f.RefreshURL)
							values += fmt.Sprintf("\n%s %s", refKey, refValue)
						}
						scopesKey := infoPageAuthenticationItemKeyColorStyle.Render("Scopes:")
						scopesValue := infoPageAuthenticationOAuthScopesStyle.Render(styledOAuthScopes(f.Scopes))
						values += fmt.Sprintf("\n%s\n%s", scopesKey, scopesValue)
						content.WriteString(infoPageAuthenticationItemStyle.Render(values))
					}
				case "openIdConnect":
					urlKey := infoPageAuthenticationItemKeyColorStyle.Render("Connect URL:")
					urlValue := infoPageAuthenticationItemValueColorStyle.Render(scheme.OpenIdConnectUrl)
					values := fmt.Sprintf("%s %s", urlKey, urlValue)
					content.WriteString(infoPageAuthenticationItemStyle.Render(values))
				}
			}
		}
	}

	content.WriteString(infoPageSeparator)

	openAPIVersion := infoPageVersionStyle.Render(fmt.Sprintf("OpenAPI Version: %s", info.OpenAPIVersion))
	content.WriteString(infoPageItemStyle.Render(openAPIVersion))

	m.viewport.SetContent(content.String())
}

func styledOAuthScopes(scopes []*topi.Scope) string {
	ss := make([]string, len(scopes))
	for i, scope := range scopes {
		scopeName := infoPageAuthenticationOAuthScopeNameColorStyle.Render(scope.Name)
		scopeDesc := infoPageAuthenticationOAuthScopeDescColorStyle.Render(scope.Detail) // fixme: render as md, consider width
		ss[i] = fmt.Sprintf("%s - %s", scopeName, scopeDesc)
	}
	return strings.Join(ss, "\n")
}

func (m *infoPageModel) selectItem(reverse bool) {
	n := infoPageSelectableNumberOfItems
	if reverse {
		m.selected = ((m.selected-1)%n + n) % n
	} else {
		m.selected = (m.selected + 1) % n
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
	case infoPageSelectableExDocsUrl:
		if m.doc.Info.ExDocsUrl == "" {
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
	case infoPageSelectableExDocsUrl:
		return openInBrowser(m.doc.Info.ExDocsUrl)
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
		m.reset()
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
