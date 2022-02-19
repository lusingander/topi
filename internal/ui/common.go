package ui

import "github.com/charmbracelet/lipgloss"

const (
	glamourTheme = "notty" // todo: fix
)

var (
	baseStyle = lipgloss.NewStyle().Margin(1, 2)

	selectedColor = lipgloss.Color("70")

	listNormalTitleColorStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#1a1a1a", Dark: "#dddddd"})

	listNormalItemStyle = lipgloss.NewStyle().
				Padding(0, 0, 0, 2)

	listNormalTitleStyle = listNormalTitleColorStyle.Copy().
				Padding(0, 0, 0, 2)

	listNormalDescColorStyle = lipgloss.NewStyle().
					Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	listNormalDescStyle = listNormalDescColorStyle.Copy().
				Padding(0, 0, 0, 2)

	listSelectedTitleColorStyle = lipgloss.NewStyle().
					Foreground(selectedColor)

	listSelectedItemStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(selectedColor).
				Padding(0, 0, 0, 1)

	listSelectedTitleStyle = listSelectedItemStyle.Copy().
				Foreground(selectedColor)

	listSelectedDescColorStyle = listSelectedTitleColorStyle.Copy().
					Foreground(selectedColor)

	listSelectedDescStyle = listSelectedItemStyle.Copy().
				Foreground(selectedColor)
)

var (
	httpMethodGetColor    = lipgloss.Color("33")
	httpMethodPostColor   = lipgloss.Color("35")
	httpMethodPutColor    = lipgloss.Color("148")
	httpMethodPatchColor  = lipgloss.Color("218")
	httpMethodDeleteColor = lipgloss.Color("172")

	httpMethodSelectedGetColor    = lipgloss.Color("31")
	httpMethodSelectedPostColor   = lipgloss.Color("29")
	httpMethodSelectedPutColor    = lipgloss.Color("112")
	httpMethodSelectedPatchColor  = lipgloss.Color("181")
	httpMethodSelectedDeleteColor = lipgloss.Color("136")

	httpMethodDeprecatedColor         = lipgloss.Color("246")
	httpMethodSelectedDeprecatedColor = lipgloss.Color("243")
)
