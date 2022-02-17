package ui

import tea "github.com/charmbracelet/bubbletea"

type selectTagMsg struct {
	tag string
}

func selectTag(tag string) tea.Cmd {
	return func() tea.Msg { return selectTagMsg{tag} }
}

type goBackMsg struct{}

func goBack() tea.Msg {
	return goBackMsg{}
}
