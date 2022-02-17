package ui

import tea "github.com/charmbracelet/bubbletea"

type selectInfoMenuMsg struct{}

func selectInfoMenu() tea.Msg {
	return selectInfoMenuMsg{}
}

type selectTagMenuMsg struct{}

func selectTagMenu() tea.Msg {
	return selectTagMenuMsg{}
}

type selectHelpMenuMsg struct{}

func selectHelpMenu() tea.Msg {
	return selectHelpMenuMsg{}
}

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
