package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type page interface{}

type menuPage struct{}

type infoPage struct{}

type tagPage struct{}

type tagApiPage struct {
	tag string
}

type helpPage struct{}

type pageStack struct {
	stack []page
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

	menuPage   menuPageModel
	infoPage   infoPageModel
	tagPage    tagPageModel
	tagApiPage tagApiPageModel
}

var _ tea.Model = (*model)(nil)

func newModel(doc *topi.Document) model {
	startPage := menuPage{}
	return model{
		doc:        doc,
		pageStack:  newPageStack(startPage),
		infoPage:   newInfoPageModel(doc),
		menuPage:   newMenuPageModel(),
		tagPage:    newTagPageModel(doc),
		tagApiPage: newTagApiPageModel(doc),
	}
}

func (m *model) SetSize(w, h int) {
	m.menuPage.SetSize(w, h)
	m.infoPage.SetSize(w, h)
	m.tagPage.SetSize(w, h)
	m.tagApiPage.SetSize(w, h)
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
		top, right, bottom, left := baseStyle.GetMargin()
		m.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	case selectInfoMenuMsg:
		m.pushPage(infoPage{})
	case selectTagMenuMsg:
		m.pushPage(tagPage{})
	case selectHelpMenuMsg:
		m.pushPage(helpPage{})
	case selectTagMsg:
		m.pushPage(tagApiPage(msg))
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
	case tagApiPage:
		m.tagApiPage, cmd = m.tagApiPage.Update(msg)
		return m, cmd
	case helpPage:
		return m, cmd // todo: impl
	default:
		return m, nil
	}
}

func (m model) View() string {
	return baseStyle.Render(m.content())
}

func (m model) content() string {
	switch m.currentPage().(type) {
	case menuPage:
		return m.menuPage.View()
	case infoPage:
		return m.infoPage.View()
	case tagPage:
		return m.tagPage.View()
	case tagApiPage:
		return m.tagApiPage.View()
	case helpPage:
		return "" // todo: impl
	default:
		return "error... :("
	}
}

func Start(doc *topi.Document) error {
	m := newModel(doc)
	p := tea.NewProgram(m, tea.WithAltScreen())
	return p.Start()
}
