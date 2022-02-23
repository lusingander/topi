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

type selectPathMenuMsg struct{}

func selectPathMenu() tea.Msg {
	return selectPathMenuMsg{}
}

type selectHelpMenuMsg struct{}

func selectHelpMenu() tea.Msg {
	return selectHelpMenuMsg{}
}

type selectAboutMenuMsg struct{}

func selectAboutMenu() tea.Msg {
	return selectAboutMenuMsg{}
}

type selectTagMsg struct {
	tag string
}

func selectTag(tag string) tea.Cmd {
	return func() tea.Msg { return selectTagMsg{tag} }
}

type selectOperationMsg struct {
	operationId string
}

func selectOperation(operationId string) tea.Cmd {
	return func() tea.Msg { return selectOperationMsg{operationId} }
}

type goBackMsg struct{}

func goBack() tea.Msg {
	return goBackMsg{}
}
