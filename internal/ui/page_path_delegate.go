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
	pathPageMethodStyle = lipgloss.NewStyle().
				Underline(true).
				Bold(true)

	pathPageMethodSelectedGetStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodSelectedGetColor)

	pathPageMethodNormalGetStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodGetColor)

	pathPageMethodSelectedPostStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodSelectedPostColor)

	pathPageMethodNormalPostStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodPostColor)

	pathPageMethodSelectedDeleteStyle = pathPageMethodStyle.Copy().
						Foreground(httpMethodSelectedDeleteColor)

	pathPageMethodNormalDeleteStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodDeleteColor)

	pathPageMethodSelectedPutStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodSelectedPutColor)

	pathPageMethodNormalPutStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodPutColor)

	pathPageMethodSelectedPatchStyle = pathPageMethodStyle.Copy().
						Foreground(httpMethodSelectedPatchColor)

	pathPageMethodNormalPatchStyle = pathPageMethodStyle.Copy().
					Foreground(httpMethodPatchColor)

	pathPageMethodDeprecatedStyle = pathPageMethodStyle.Copy().
					Strikethrough(true)

	pathPageMethodSelectedDeprecatedStyle = pathPageMethodDeprecatedStyle.Copy().
						Foreground(httpMethodSelectedDeprecatedColor)

	pathPageMethodNormalDeprecatedStyle = pathPageMethodDeprecatedStyle.Copy().
						Foreground(httpMethodDeprecatedColor)

	pathPagePathDeprecatedStyle = lipgloss.NewStyle().
					Strikethrough(true)

	pathPagePathSelectedDeprecatedStyle = pathPagePathDeprecatedStyle.Copy().
						Foreground(selectedColor)

	pathPagePathNormalDeprecatedStyle = pathPagePathDeprecatedStyle.Copy().
						Foreground(lipgloss.Color("246"))
)

type pathPageListItem struct {
	path *topi.Path
}

var _ list.Item = (*pathPageListItem)(nil)

func (i pathPageListItem) FilterValue() string {
	return i.path.UriPath
}

func (i pathPageListItem) styledTitle(selected bool) string {
	var method, path string
	if i.path.Deprecated {
		if selected {
			method = pathPageMethodSelectedDeprecatedStyle.Render(i.path.Method)
			path = pathPagePathSelectedDeprecatedStyle.Render(i.path.UriPath)
		} else {
			method = pathPageMethodNormalDeprecatedStyle.Render(i.path.Method)
			path = pathPagePathNormalDeprecatedStyle.Render(i.path.UriPath)
		}
		return fmt.Sprintf("%s %s", padding.String(method, 7), path)
	}
	switch i.path.Method {
	case http.MethodGet:
		if selected {
			method = pathPageMethodSelectedGetStyle.Render(i.path.Method)
		} else {
			method = pathPageMethodNormalGetStyle.Render(i.path.Method)
		}
	case http.MethodPost:
		if selected {
			method = pathPageMethodSelectedPostStyle.Render(i.path.Method)
		} else {
			method = pathPageMethodNormalPostStyle.Render(i.path.Method)
		}
	case http.MethodDelete:
		if selected {
			method = pathPageMethodSelectedDeleteStyle.Render(i.path.Method)
		} else {
			method = pathPageMethodNormalDeleteStyle.Render(i.path.Method)
		}
	case http.MethodPut:
		if selected {
			method = pathPageMethodSelectedPutStyle.Render(i.path.Method)
		} else {
			method = pathPageMethodNormalPutStyle.Render(i.path.Method)
		}
	case http.MethodPatch:
		if selected {
			method = pathPageMethodSelectedPatchStyle.Render(i.path.Method)
		} else {
			method = pathPageMethodNormalPatchStyle.Render(i.path.Method)
		}
	default:
		method = pathPageMethodStyle.Render(i.path.Method)
	}
	if selected {
		path = listSelectedTitleColorStyle.Render(i.path.UriPath)
	} else {
		path = listNormalTitleColorStyle.Render(i.path.UriPath)
	}
	return fmt.Sprintf("%s %s", padding.String(method, 7), path)
}

func (i pathPageListItem) styledDesc(selected bool, width int) string {
	desc := truncateWithTail(i.path.Summary, uint(width))
	if selected {
		desc = listSelectedDescColorStyle.Render(desc)
	} else {
		desc = listNormalDescColorStyle.Render(desc)
	}
	return desc
}

type pathPageListDelegate struct{}

var _ list.ItemDelegate = (*pathPageListDelegate)(nil)

func newPathPageListDelegate() pathPageListDelegate {
	return pathPageListDelegate{}
}

func (d pathPageListDelegate) Height() int {
	return 2
}

func (d pathPageListDelegate) Spacing() int {
	return 1
}

func (d pathPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d pathPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(pathPageListItem)
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
