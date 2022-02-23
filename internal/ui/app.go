package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
)

var (
	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("65")).
			Padding(0, 1)

	footerStyle = lipgloss.NewStyle()

	statusbarFileNameStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("65")).
				Padding(0, 1)

	statusbarSpaceColorStyle = lipgloss.NewStyle().
					Background(lipgloss.Color("237"))

	statusbarLowerStyle = lipgloss.NewStyle()
)

type page interface {
	crumb() string
}

type menuPage struct{}

func (menuPage) crumb() string { return topi.AppName }

type infoPage struct{}

func (infoPage) crumb() string { return "info" }

type tagPage struct{}

func (tagPage) crumb() string { return "tags" }

type tagPathsPage struct {
	tag string
}

func (p tagPathsPage) crumb() string { return p.tag }

type operationPage struct {
	operationId string
}

func (p operationPage) crumb() string { return p.operationId } // fixme

type helpPage struct{}

func (helpPage) crumb() string { return "help" }

type aboutPage struct{}

func (aboutPage) crumb() string { return "about" }

type pageStack struct {
	stack []page
}

func (s pageStack) crumbs() []string {
	ret := make([]string, len(s.stack))
	for i, p := range s.stack {
		ret[i] = p.crumb()
	}
	return ret
}

func newPageStack(p page) *pageStack {
	return &pageStack{
		stack: []page{p},
	}
}

func (s *pageStack) pushPage(p page) {
	s.stack = append(s.stack, p)
}

func (s *pageStack) popPage() page {
	l := len(s.stack)
	if l <= 1 {
		return nil
	}
	p := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return p
}

func (s *pageStack) currentPage() page {
	return s.stack[len(s.stack)-1]
}

type model struct {
	doc *topi.Document

	*pageStack

	menuPage      menuPageModel
	infoPage      infoPageModel
	tagPage       tagPageModel
	tagPathsPage  tagPathsPageModel
	operationPage operationPageModel
	helpPage      helpPageModel
	aboutPage     aboutPageModel

	width, height int
}

var _ tea.Model = (*model)(nil)

func newModel(doc *topi.Document) model {
	startPage := menuPage{}
	return model{
		doc:           doc,
		pageStack:     newPageStack(startPage),
		infoPage:      newInfoPageModel(doc),
		menuPage:      newMenuPageModel(),
		tagPage:       newTagPageModel(doc),
		tagPathsPage:  newTagPathsPageModel(doc),
		operationPage: newOperationPageModel(doc),
		helpPage:      newHelpPageModel(),
		aboutPage:     newAboutPageModel(),
	}
}

func (m *model) SetSize(w, h int) {
	m.width, m.height = w, h

	t, r, b, l := baseStyle.GetMargin()
	w = w - r - l
	h = h - t - b
	h = h - 3 // :(

	m.menuPage.SetSize(w, h)
	m.infoPage.SetSize(w, h)
	m.tagPage.SetSize(w, h)
	m.tagPathsPage.SetSize(w, h)
	m.operationPage.SetSize(w, h)
	m.helpPage.SetSize(w, h)
	m.aboutPage.SetSize(w, h)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	case selectInfoMenuMsg:
		m.pushPage(infoPage{})
	case selectTagMenuMsg:
		m.pushPage(tagPage{})
	case selectHelpMenuMsg:
		m.pushPage(helpPage{})
	case selectAboutMenuMsg:
		m.pushPage(aboutPage{})
	case selectTagMsg:
		m.pushPage(tagPathsPage(msg))
	case selectOperationMsg:
		m.pushPage(operationPage(msg))
	case goBackMsg:
		m.popPage()
	}
	switch m.currentPage().(type) {
	case menuPage:
		m.menuPage, cmd = m.menuPage.Update(msg)
		return m, cmd
	case infoPage:
		m.infoPage, cmd = m.infoPage.Update(msg)
		return m, cmd
	case tagPage:
		m.tagPage, cmd = m.tagPage.Update(msg)
		return m, cmd
	case tagPathsPage:
		m.tagPathsPage, cmd = m.tagPathsPage.Update(msg)
		return m, cmd
	case operationPage:
		m.operationPage, cmd = m.operationPage.Update(msg)
		return m, cmd
	case helpPage:
		m.helpPage, cmd = m.helpPage.Update(msg)
		return m, cmd
	case aboutPage:
		m.aboutPage, cmd = m.aboutPage.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	header := m.appHeader()
	content := baseStyle.Render(m.content())
	footer := m.appFooter()
	return lipgloss.JoinVertical(lipgloss.Top, header, content, footer)
}

func (m model) content() string {
	switch m.currentPage().(type) {
	case menuPage:
		return m.menuPage.View()
	case infoPage:
		return m.infoPage.View()
	case tagPage:
		return m.tagPage.View()
	case tagPathsPage:
		return m.tagPathsPage.View()
	case operationPage:
		return m.operationPage.View()
	case helpPage:
		return m.helpPage.View()
	case aboutPage:
		return m.aboutPage.View()
	default:
		return "error... :("
	}
}

func (m model) appHeader() string {
	bd := strings.Join(m.crumbs(), " > ")
	return headerStyle.Render(bd)
}

func (m model) appFooter() string {
	w := m.width
	if w == 0 {
		return ""
	}
	name := statusbarFileNameStyle.Render(m.doc.Meta.FileName)
	spaces := strings.Repeat(" ", w-lipgloss.Width(name))
	spaces = statusbarSpaceColorStyle.Render(spaces)
	u := name + spaces
	status := m.statusString()
	l := statusbarLowerStyle.Render(status)
	return footerStyle.Render(u + "\n" + l)
}

func (m model) statusString() string {
	switch m.currentPage().(type) {
	case menuPage:
		return ""
	case infoPage:
		return ""
	case tagPage:
		return m.tagPage.statusString()
	case tagPathsPage:
		return m.tagPathsPage.statusString()
	case operationPage:
		return ""
	case helpPage:
		return ""
	case aboutPage:
		return ""
	default:
		return "error... :("
	}
}

func Start(doc *topi.Document) error {
	m := newModel(doc)
	p := tea.NewProgram(m, tea.WithAltScreen())
	return p.Start()
}
