package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
)

var (
	operationPageMethodStyle = lipgloss.NewStyle().
					Bold(true)

	operationPageMethodGetStyle = operationPageMethodStyle.Copy().
					Foreground(httpMethodGetColor)

	operationPageMethodPostStyle = operationPageMethodStyle.Copy().
					Foreground(httpMethodPostColor)

	operationPageMethodPutStyle = operationPageMethodStyle.Copy().
					Foreground(httpMethodPutColor)

	operationPageMethodPatchStyle = operationPageMethodStyle.Copy().
					Foreground(httpMethodPatchColor)

	operationPageMethodDeleteStyle = operationPageMethodStyle.Copy().
					Foreground(httpMethodDeleteColor)

	operationPageMethodDeprecatedStyle = operationPageMethodStyle.Copy().
						Foreground(httpMethodDeprecatedColor)

	operationPageDeprecatedMarkerStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("208")).
						Bold(true).
						Margin(0, 0, 0, 2)

	operationPageSectionHeaderStyle = lipgloss.NewStyle().
					Foreground(lipgloss.Color("70")).
					Underline(true)

	operationPageSectionSubHeaderStyle = operationPageSectionHeaderStyle.Copy().
						Margin(0, 0, 0, 1)

	operationPageParameterItemsStyle = operationPageItemStyle.Copy().
						Margin(0, 0, 0, 2)

	operationPageParameterRequiredMarkerColorStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("168"))

	operationPageParameterTypeColorStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("246"))

	operationPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

type operationPageModel struct {
	doc           *topi.Document
	operation     *topi.Path
	viewport      viewport.Model
	delegateKeys  operationPageDelegateKeyMap
	width, height int
}

func newOperationPageModel(doc *topi.Document) operationPageModel {
	m := operationPageModel{
		doc:       doc,
		operation: nil,
	}
	m.delegateKeys = newOperationPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	return m
}

type operationPageDelegateKeyMap struct {
	back key.Binding
}

func newOperationPageDelegateKeyMap() operationPageDelegateKeyMap {
	return operationPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *operationPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.viewport.Width, m.viewport.Height = w, h
	m.updateContent()
}

func (m *operationPageModel) updateOperation(operationId string) {
	m.operation = m.doc.FindPathByOperationId(operationId)
}

func (m *operationPageModel) updateContent() {
	op := m.operation
	if op == nil {
		return
	}

	glamourWidth := m.width - 10
	r, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle(glamourTheme),
		glamour.WithWordWrap(glamourWidth),
	)

	var content strings.Builder

	method := m.styledMethod()
	path := op.UriPath
	mp := fmt.Sprintf("%s %s", method, path)
	if op.Deprecated {
		mp += operationPageDeprecatedMarkerStyle.Render("Deprecated")
	}
	content.WriteString(operationPageItemStyle.Render(mp))

	if op.Summary != "" {
		summary := op.Summary
		content.WriteString(operationPageItemStyle.Render(summary))
	}

	if op.Description != "" {
		desc, _ := r.Render(op.Description)
		desc = operationPageItemStyle.Render(desc)
		content.WriteString(desc)
	}

	requestSectionHeader := operationPageSectionHeaderStyle.Render("Request")
	content.WriteString(operationPageItemStyle.Render(requestSectionHeader))

	if len(op.PathParameters) > 0 {
		pathParamSectionHeader := operationPageSectionSubHeaderStyle.Render("Path parameters")
		content.WriteString(operationPageItemStyle.Render(pathParamSectionHeader))
		content.WriteString(operationPageParameterItemsStyle.Render(m.styledParams(op.PathParameters)))
	}

	if len(op.QueryParameters) > 0 {
		queryParamSectionHeader := operationPageSectionSubHeaderStyle.Render("Query parameters")
		content.WriteString(operationPageItemStyle.Render(queryParamSectionHeader))
		content.WriteString(operationPageParameterItemsStyle.Render(m.styledParams(op.QueryParameters)))
	}

	if len(op.HeaderParameters) > 0 {
		headerParamSectionHeader := operationPageSectionSubHeaderStyle.Render("Header parameters")
		content.WriteString(operationPageItemStyle.Render(headerParamSectionHeader))
		content.WriteString(operationPageParameterItemsStyle.Render(m.styledParams(op.HeaderParameters)))
	}

	if len(op.CookieParameters) > 0 {
		cookieParamSectionHeader := operationPageSectionSubHeaderStyle.Render("Cookie parameters")
		content.WriteString(operationPageItemStyle.Render(cookieParamSectionHeader))
		content.WriteString(operationPageParameterItemsStyle.Render(m.styledParams(op.CookieParameters)))
	}

	responseSectionHeader := operationPageSectionHeaderStyle.Render("Response")
	content.WriteString(operationPageItemStyle.Render(responseSectionHeader))

	m.viewport.SetContent(content.String())
}

func (operationPageModel) styledParams(params []*topi.Parameter) string {
	strs := make([]string, len(params))
	for i, param := range params {
		p := param.Name
		if param.Required {
			p += operationPageParameterRequiredMarkerColorStyle.Render("*")
		}
		if param.Schema != nil {
			p += "  "
			if param.Schema.Type != "" {
				p += operationPageParameterTypeColorStyle.Render(param.Schema.Type)
			}
			if param.Schema.Format != "" {
				p += operationPageParameterTypeColorStyle.Render(fmt.Sprintf("(%s)", param.Schema.Type))
			}
		}
		strs[i] = p
	}
	return strings.Join(strs, "\n")
}

func (m operationPageModel) styledMethod() string {
	method := m.operation.Method
	if m.operation.Deprecated {
		return operationPageMethodDeprecatedStyle.Render(method)
	}
	switch method {
	case http.MethodGet:
		return operationPageMethodGetStyle.Render(method)
	case http.MethodPost:
		return operationPageMethodPostStyle.Render(method)
	case http.MethodPut:
		return operationPageMethodPutStyle.Render(method)
	case http.MethodPatch:
		return operationPageMethodPatchStyle.Render(method)
	case http.MethodDelete:
		return operationPageMethodDeleteStyle.Render(method)
	default:
		return operationPageMethodStyle.Render(method)
	}
}

func (m operationPageModel) Init() tea.Cmd {
	return nil
}

func (m operationPageModel) Update(msg tea.Msg) (operationPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		}
	case selectOperationMsg:
		m.updateOperation(msg.operationId)
		m.updateContent()
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m operationPageModel) View() string {
	return m.viewport.View()
}
