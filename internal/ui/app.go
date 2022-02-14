package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/lusingander/topi/internal/topi"
)

type page int

const (
	tagPage page = iota
	tagApiPage
)

type model struct {
	doc *topi.Document

	currentPage page

	tagPage    tagPageModel
	tagApiPage tagApiPageModel
}

var _ tea.Model = (*model)(nil)

func newModel(doc *topi.Document) model {
	return model{
		doc:         doc,
		currentPage: tagPage,
		tagPage:     newTagPageModel(doc),
		tagApiPage:  newTagApiPageModel(doc),
	}
}

func (m *model) SetSize(w, h int) {
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
	case selectTagMsg:
		m.currentPage = tagApiPage
	case goBackTagPageMsg:
		m.currentPage = tagPage
	}
	switch m.currentPage {
	case tagPage:
		m.tagPage, cmd = m.tagPage.Update(msg)
		return m, cmd
	case tagApiPage:
		m.tagApiPage, cmd = m.tagApiPage.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m model) View() string {
	return baseStyle.Render(m.content())
}

func (m model) content() string {
	switch m.currentPage {
	case tagPage:
		return m.tagPage.View()
	case tagApiPage:
		return m.tagApiPage.View()
	default:
		return "error... :("
	}
}

func Start(doc *topi.Document) error {
	m := newModel(doc)
	p := tea.NewProgram(m, tea.WithAltScreen())
	return p.Start()
}
