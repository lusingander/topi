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

	opearationPageRequestBodyMediaTypeColorStyle = lipgloss.NewStyle().
							Foreground(lipgloss.Color("70"))

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

	operationPageSectionSubHeaderSuccessStatusCodeStyle = operationPageSectionSubHeaderStyle.Copy().
								Foreground(lipgloss.Color("77"))

	operationPageSectionSubHeaderErrorStatusCodeStyle = operationPageSectionSubHeaderStyle.Copy().
								Foreground(lipgloss.Color("168"))

	operationPageSectionSubHeaderDefaultStatusCodeStyle = operationPageSectionSubHeaderStyle.Copy().
								Foreground(lipgloss.Color("32"))

	operationPageSecuritySchemeScopeStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("203")).
						Background(lipgloss.Color("236"))

	operationPageSchemaIndentColorStyle = lipgloss.NewStyle().
						Foreground(lipgloss.Color("238"))

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
	// fixme: operationId is not required field...
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

	if len(op.Security) > 0 {
		requestSectionHeader := operationPageSectionHeaderStyle.Render("Security Requirements")
		content.WriteString(operationPageItemStyle.Render(requestSectionHeader))

		rs := styledSecurityRequirements(op.Security)
		content.WriteString(operationPageParameterItemsStyle.Render(rs))
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

	if op.RequestBody != nil && len(op.RequestBody.Conetnt) > 0 {
		for _, c := range op.RequestBody.Conetnt {
			requestBodyMediaType := opearationPageRequestBodyMediaTypeColorStyle.Render(fmt.Sprintf("[%s]", c.MediaType))
			requestBodySectionHeader := operationPageSectionSubHeaderStyle.Render("Request body")
			content.WriteString(operationPageItemStyle.Render(fmt.Sprintf("%s  %s", requestBodySectionHeader, requestBodyMediaType)))
			content.WriteString(operationPageParameterItemsStyle.Render(styledSchema(c.Schema, 0, false)))
		}
	}

	responseSectionHeader := operationPageSectionHeaderStyle.Render("Response")
	content.WriteString(operationPageItemStyle.Render(responseSectionHeader))

	for _, response := range op.Responses {
		var statusCode string
		if response.Success() {
			statusCode = operationPageSectionSubHeaderSuccessStatusCodeStyle.Render(response.StatusCode)
		} else if response.Error() {
			statusCode = operationPageSectionSubHeaderErrorStatusCodeStyle.Render(response.StatusCode)
		} else {
			statusCode = operationPageSectionSubHeaderDefaultStatusCodeStyle.Render(response.StatusCode)
		}

		desc := ""
		if response.Description != "" {
			desc = response.Description // fixme: render as md, consider width
		}
		content.WriteString(operationPageItemStyle.Render(fmt.Sprintf("%s  %s", statusCode, desc)))

		if len(response.Headers) > 0 {
			requestHeadersHeader := operationPageSectionSubHeaderStyle.Render("Response headers")
			content.WriteString(operationPageItemStyle.Render(requestHeadersHeader))
			content.WriteString(operationPageParameterItemsStyle.Render(m.styledHeaders(response.Headers)))
		}

		for _, c := range response.Conetnt {
			if c.Schema == nil {
				continue
			}
			requestBodyMediaTypeHeader := operationPageSectionSubHeaderStyle.Render("Response schema")
			requestBodyMediaType := opearationPageRequestBodyMediaTypeColorStyle.Render(fmt.Sprintf("[%s]", c.MediaType))
			content.WriteString(operationPageItemStyle.Render(fmt.Sprintf("%s  %s", requestBodyMediaTypeHeader, requestBodyMediaType)))
			content.WriteString(operationPageParameterItemsStyle.Render(styledSchema(c.Schema, 0, true)))
		}
	}

	m.viewport.SetContent(content.String())
}

func styledSecurityRequirements(requirements []*topi.SecurityRequirement) string {
	ss := make([]string, len(requirements))
	for i, requirement := range requirements {
		schemes := make([]string, len(requirement.Schemes))
		for j, scheme := range requirement.Schemes {
			schemes[j] = styledSecurityScheme(scheme)
		}
		ss[i] = strings.Join(schemes, " + ")
	}
	return strings.Join(ss, "\n or\n")
}

func styledSecurityScheme(sc *topi.SecurityRequirementScheme) string {
	s := sc.Key
	if len(sc.Scopes) > 0 {
		scopes := make([]string, len(sc.Scopes))
		for i, scope := range sc.Scopes {
			scopes[i] = operationPageSecuritySchemeScopeStyle.Render(scope)
		}
		s += fmt.Sprintf(" (%s)", strings.Join(scopes, ", "))
	}
	return s
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
	nameAreaWidth += 2 // requred marker + buf

	for _, param := range params {
		ss := styledSingleParam(param.Schema, param.Name, param.Description, param.Required, param.Deprecated, nameAreaWidth, 0)
		strs = append(strs, ss...)
	}
	return strings.Join(strs, "\n")
}

func (operationPageModel) styledHeaders(headers []*topi.Header) string {
	strs := make([]string, 0)

	nameAreaWidth := 0
	for _, header := range headers {
		w := len(header.Name)
		if nameAreaWidth < w {
			nameAreaWidth = w
		}
	}
	nameAreaWidth += 2 // requred marker + buf

	for _, header := range headers {
		ss := styledSingleParam(header.Parameter.Schema, header.Name, header.Parameter.Description, header.Parameter.Required, header.Parameter.Deprecated, nameAreaWidth, 0)
		strs = append(strs, ss...)
	}
	return strings.Join(strs, "\n")
}

func styledSchema(sc *topi.Schema, indentLevel int, read bool) string {
	if sc.Type == "object" {

		nameAreaWidth := 0
		for name := range sc.Properties {
			w := len(name)
			if nameAreaWidth < w {
				nameAreaWidth = w
			}
		}
		nameAreaWidth += 2 // requred marker + buf

		strs := make([]string, 0)
		for name, prop := range sc.Properties { // fixme: fix order
			if read {
				if prop.WriteOnly {
					continue
				}
			} else {
				if prop.ReadOnly {
					continue
				}
			}
			required := containsString(name, sc.Required)
			ss := styledSingleParam(prop, name, prop.Description, required, prop.Deprecated, nameAreaWidth, indentLevel)
			strs = append(strs, ss...)

			if prop.Type == "object" {
				ss := styledProperties(prop, indentLevel+1, read)
				strs = append(strs, ss...)
			}
			if prop.Type == "array" && prop.Items.Type == "object" {
				ss := styledProperties(prop.Items, indentLevel+1, read)
				strs = append(strs, ss...)
			}
		}
		return strings.Join(strs, "\n")
	}
	if sc.Type == "array" && sc.Items.Type == "object" {
		s := schemaTypeString(sc)
		t := styledSchema(sc.Items, indentLevel+1, read)
		return strings.Join([]string{s, t}, "\n")
	}
	return schemaTypeString(sc)
}

func styledProperties(sc *topi.Schema, indentLevel int, read bool) []string {
	strs := make([]string, 0)
	props := sc.Properties

	nameAreaWidth := 0
	for name := range props {
		w := len(name)
		if nameAreaWidth < w {
			nameAreaWidth = w
		}
	}
	nameAreaWidth += 2 // requred marker + buf

	for name, prop := range props { // fixme: fix order
		if read {
			if prop.WriteOnly {
				continue
			}
		} else {
			if prop.ReadOnly {
				continue
			}
		}
		required := containsString(name, sc.Required)
		ss := styledSingleParam(prop, name, prop.Description, required, prop.Deprecated, nameAreaWidth, indentLevel)
		strs = append(strs, ss...)

		if prop.Type == "object" {
			ss := styledProperties(prop, indentLevel+1, read)
			strs = append(strs, ss...)
		}
		if prop.Type == "array" && prop.Items.Type == "object" {
			ss := styledProperties(prop.Items, indentLevel+1, read)
			strs = append(strs, ss...)
		}
	}

	return strs
}

func styledSingleParam(schema *topi.Schema, name, description string, required, deprecated bool, nameAreaWidth, indentLevel int) []string {
	strs := make([]string, 0)

	schemaIndent := strings.Repeat(">>", indentLevel)
	descIndent := strings.Repeat(" ", nameAreaWidth+len(schemaIndent))

	var s strings.Builder

	s.WriteString(operationPageSchemaIndentColorStyle.Render(schemaIndent))

	if deprecated {
		name = operationPageParameterDeprecatedNameStyle.Render(name)
	}
	if required {
		name += operationPageParameterRequiredMarkerColorStyle.Render("*")
	}
	s.WriteString(padding.String(name, uint(nameAreaWidth)))

	if schema != nil {
		schemaType := schemaTypeString(schema)
		s.WriteString(operationPageParameterTypeColorStyle.Render(schemaType))
	}

	if deprecated {
		s.WriteString(" ")
		s.WriteString(operationPageDeprecatedMarkerStyle.Render("Deprecated"))
	}

	strs = append(strs, s.String())

	if description != "" {
		var s strings.Builder
		s.WriteString(descIndent)
		s.WriteString(description) // fixme: render as md, consider width
		strs = append(strs, s.String())
	}

	if schema != nil {
		if schema.Default != nil {
			var s strings.Builder
			k := operationPageParameterPropertyKeyStyle.Render("Default:")
			v := operationPageParameterPropertyValueStyle.Render(fmt.Sprintf("%v", schema.Default))
			s.WriteString(descIndent)
			s.WriteString(fmt.Sprintf("%s %s", k, v))
			strs = append(strs, s.String())
		}

		if len(schema.Enum) > 0 {
			var s strings.Builder
			k := operationPageParameterPropertyKeyStyle.Render("Enum:")
			v := operationPageParameterPropertyValueStyle.Render(sliceString(schema.Enum))
			s.WriteString(descIndent)
			s.WriteString(fmt.Sprintf("%s %s", k, v))
			strs = append(strs, s.String())
		}

		if schema.Type == "array" && len(schema.Items.Enum) > 0 {
			var s strings.Builder
			k := operationPageParameterPropertyKeyStyle.Render("Items Enum:")
			v := operationPageParameterPropertyValueStyle.Render(sliceString(schema.Items.Enum))
			s.WriteString(descIndent)
			s.WriteString(fmt.Sprintf("%s %s", k, v))
			strs = append(strs, s.String())
		}

		constraints := schemaConstraintStrings(schema)
		if len(constraints) > 0 {
			var s strings.Builder
			k := operationPageParameterPropertyKeyStyle.Render("Constraints:")
			v := operationPageParameterPropertyValueStyle.Render(strings.Join(constraints, ", "))
			s.WriteString(descIndent)
			s.WriteString(fmt.Sprintf("%s %s", k, v))
			strs = append(strs, s.String())
		}
	}
	return strs
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
