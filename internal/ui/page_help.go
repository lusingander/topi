package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

const helpPageContent = `# Help page

*press j to scroll down, and k to scroll up*

*press backspace or ? to go back to the previous page*

## Keybindings

### Common

common keybindings for all pages

|Key|Description|
|-|-|
|Backspace|back to perv page|
|Ctrl+c|quit|
|?|show help page (this page)|

### List page

keybindings for list-syle pages 

|Key|Description|
|-|-|
|j|cursor down|
|k|cursor up|
|f l|next page|
|b h|prev page|
|g|go to start|
|G|go to end|
|/|Enter filtering mode|
|Enter|(default) select item, (filtering) apply filter|
|Esc|(filtering) cancel filter, (filter applied) remove filter|

## Document page

keybindings for document pages 

|Key|Description|
|-|-|
|j|page down one line|
|k|page up one line|
|f|page down|
|b|page up|
|d|half page down|
|u|half page up|
|Tab|select link|
|x|open selecting link|
`

type helpPageModel struct {
	viewport      viewport.Model
	delegateKeys  helpPageDelegateKeyMap
	width, height int
}

func newHelpPageModel() helpPageModel {
	m := helpPageModel{}
	m.delegateKeys = newHelpPageDelegateKeyMap()
	m.viewport = viewport.New(0, 0)
	return m
}

type helpPageDelegateKeyMap struct {
	back key.Binding
}

func newHelpPageDelegateKeyMap() helpPageDelegateKeyMap {
	return helpPageDelegateKeyMap{
		back: key.NewBinding(
			key.WithKeys("backspace", "ctrl+h"),
			key.WithHelp("backspace", "back"),
		),
	}
}

func (m *helpPageModel) SetSize(w, h int) {
	m.width, m.height = w, h
	m.viewport.Width, m.viewport.Height = w, h
	m.updateContent()
}

func (m *helpPageModel) reset() {
	m.viewport.GotoTop()
}

func (m *helpPageModel) updateContent() {
	r, _ := markdownRenderer(m.width - 10)
	content, _ := r.Render(helpPageContent)
	m.viewport.SetContent(content)
}

func (m helpPageModel) Init() tea.Cmd {
	return nil
}

func (m helpPageModel) Update(msg tea.Msg) (helpPageModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.delegateKeys.back):
			return m, goBack
		}
	case selectHelpHelpMenuMsg, toggleHelpMsg:
		m.reset()
		m.updateContent()
		return m, nil
	}
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m helpPageModel) View() string {
	return m.viewport.View()
}
