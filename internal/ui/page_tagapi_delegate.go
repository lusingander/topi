package ui

import (
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lusingander/topi/internal/topi"
	"github.com/muesli/reflow/padding"
)

var (
	tagApiPageMethodStyle = lipgloss.NewStyle().
				Underline(true).
				Bold(true)

	tagApiPageMethodSelectedGetStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("31"))

	tagApiPageMethodNormalGetStyle = tagApiPageMethodStyle.Copy().
					Foreground(lipgloss.Color("33"))

	tagApiPageMethodSelectedPostStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("29"))

	tagApiPageMethodNormalPostStyle = tagApiPageMethodStyle.Copy().
					Foreground(lipgloss.Color("35"))

	tagApiPageMethodSelectedDeleteStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("136"))

	tagApiPageMethodNormalDeleteStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("172"))

	tagApiPageMethodSelectedPutStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("112"))

	tagApiPageMethodNormalPutStyle = tagApiPageMethodStyle.Copy().
					Foreground(lipgloss.Color("148"))

	tagApiPageMethodSelectedPatchStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("181"))

	tagApiPageMethodNormalPatchStyle = tagApiPageMethodStyle.Copy().
						Foreground(lipgloss.Color("218"))
)

type tagApiPageListItem struct {
	path *topi.Path
}

var _ list.Item = (*tagApiPageListItem)(nil)

func (i tagApiPageListItem) FilterValue() string {
	return i.path.UriPath
}

func (i tagApiPageListItem) styledTitle(selected bool) string {
	var method, path string
	switch i.path.Method {
	case http.MethodGet:
		if selected {
			method = tagApiPageMethodSelectedGetStyle.Render(i.path.Method)
		} else {
			method = tagApiPageMethodNormalGetStyle.Render(i.path.Method)
		}
	case http.MethodPost:
		if selected {
			method = tagApiPageMethodSelectedPostStyle.Render(i.path.Method)
		} else {
			method = tagApiPageMethodNormalPostStyle.Render(i.path.Method)
		}
	case http.MethodDelete:
		if selected {
			method = tagApiPageMethodSelectedDeleteStyle.Render(i.path.Method)
		} else {
			method = tagApiPageMethodNormalDeleteStyle.Render(i.path.Method)
		}
	case http.MethodPut:
		if selected {
			method = tagApiPageMethodSelectedPutStyle.Render(i.path.Method)
		} else {
			method = tagApiPageMethodNormalPutStyle.Render(i.path.Method)
		}
	case http.MethodPatch:
		if selected {
			method = tagApiPageMethodSelectedPatchStyle.Render(i.path.Method)
		} else {
			method = tagApiPageMethodNormalPatchStyle.Render(i.path.Method)
		}
	default:
		method = tagApiPageMethodStyle.Render(i.path.Method)
	}
	if selected {
		path = listSelectedTitleColorStyle.Render(i.path.UriPath)
	} else {
		path = listNormalTitleColorStyle.Render(i.path.UriPath)
	}
	return fmt.Sprintf("%s %s", padding.String(method, 7), path)
}

func (i tagApiPageListItem) styledDesc(selected bool, width int) string {
	desc := truncateWithTail(i.path.Summary, uint(width))
	if selected {
		desc = listSelectedDescColorStyle.Render(desc)
	} else {
		desc = listNormalDescColorStyle.Render(desc)
	}
	return desc
}

type tagApiPageListDelegate struct{}

var _ list.ItemDelegate = (*tagApiPageListDelegate)(nil)

func newTagApiPageListDelegate() tagApiPageListDelegate {
	return tagApiPageListDelegate{}
}

func (d tagApiPageListDelegate) Height() int {
	return 2
}

func (d tagApiPageListDelegate) Spacing() int {
	return 1
}

func (d tagApiPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d tagApiPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(tagApiPageListItem)
	selected := index == m.Index()

	width := m.Width() - listNormalTitleStyle.GetPaddingLeft() - listNormalTitleStyle.GetPaddingRight()

	title := i.styledTitle(selected)
	desc := i.styledDesc(selected, width)

	if selected {
		title = listSelectedItemStyle.Render(title)
		desc = listSelectedItemStyle.Render(desc)
	} else {
		title = listNormalItemStyle.Render(title)
		desc = listNormalItemStyle.Render(desc)
	}

	fmt.Fprintf(w, "%s\n%s", title, desc)
}
