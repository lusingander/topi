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
						Bold(true)

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

	var content strings.Builder

	method := m.styledMethod()
	path := op.UriPath
	mp := fmt.Sprintf("%s %s", method, path)
	if op.Deprecated {
		mp += operationPageDeprecatedMarkerStyle.Render(" Deprecated")
	}
	content.WriteString(operationPageItemStyle.Render(mp))

	summary := op.Summary
	content.WriteString(operationPageItemStyle.Render(summary))

	m.viewport.SetContent(content.String())
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
