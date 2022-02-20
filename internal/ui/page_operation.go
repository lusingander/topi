package ui

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
	"github.com/muesli/reflow/padding"
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

	operationPageParameterDeprecatedNameStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("246")).
							Strikethrough(true)

	operationPageParameterPropertyKeyStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("143"))

	operationPageParameterPropertyValueStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("167"))

	operationPageItemStyle = lipgloss.NewStyle().
				Padding(1, 2)
)

var (
	operationPageSeparator = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(1, 2).
		Render("----------")
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

func (m *operationPageModel) reset() {
	m.viewport.GotoTop()
}

func (m *operationPageModel) updateOperation(operationId string) {
	m.operation = m.doc.FindPathByOperationId(operationId)
}

func (m *operationPageModel) updateContent() {
	op := m.operation
	if op == nil {
		return
	}

	r, _ := markdownRenderer(m.width - 10)

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
	content.WriteString(operationPageSeparator)

	if op.Description != "" {
		desc, _ := r.Render(op.Description)
		desc = operationPageItemStyle.Render(desc)
		content.WriteString(desc)

		content.WriteString(operationPageSeparator)
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
	strs := make([]string, 0)

	nameAreaWidth := 0
	for _, param := range params {
		w := len(param.Name)
		if nameAreaWidth < w {
			nameAreaWidth = w
		}
	}
	nameAreaWidth += 1 // requred marker
	schemaIndent := "  "
	descIndent := strings.Repeat(" ", nameAreaWidth+len(schemaIndent))

	for _, param := range params {
		var s strings.Builder

		name := param.Name
		if param.Deprecated {
			name = operationPageParameterDeprecatedNameStyle.Render(name)
		}
		if param.Required {
			name += operationPageParameterRequiredMarkerColorStyle.Render("*")
		}
		s.WriteString(padding.String(name, uint(nameAreaWidth)))

		if param.Schema != nil {
			schemaType := schemaTypeString(param.Schema)
			s.WriteString(schemaIndent)
			s.WriteString(operationPageParameterTypeColorStyle.Render(schemaType))
		}

		if param.Deprecated {
			s.WriteString(" ")
			s.WriteString(operationPageDeprecatedMarkerStyle.Render("Deprecated"))
		}

		strs = append(strs, s.String())

		if param.Description != "" {
			var s strings.Builder
			s.WriteString(descIndent)
			s.WriteString(param.Description) // fixme: render as md, consider width
			strs = append(strs, s.String())
		}

		if param.Schema != nil {
			if param.Schema.Default != nil {
				var s strings.Builder
				k := operationPageParameterPropertyKeyStyle.Render("Default:")
				v := operationPageParameterPropertyValueStyle.Render(fmt.Sprintf("%v", param.Schema.Default))
				s.WriteString(descIndent)
				s.WriteString(fmt.Sprintf("%s %s", k, v))
				strs = append(strs, s.String())
			}

			constraints := schemaConstraintStrings(param.Schema)
			if len(constraints) > 0 {
				var s strings.Builder
				k := operationPageParameterPropertyKeyStyle.Render("Constraints:")
				v := operationPageParameterPropertyValueStyle.Render(strings.Join(constraints, ", "))
				s.WriteString(descIndent)
				s.WriteString(fmt.Sprintf("%s %s", k, v))
				strs = append(strs, s.String())
			}
		}
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
		m.reset()
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
