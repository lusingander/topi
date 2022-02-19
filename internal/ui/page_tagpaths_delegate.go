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
	tagPathsPageMethodStyle = lipgloss.NewStyle().
				Underline(true).
				Bold(true)

	tagPathsPageMethodSelectedGetStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodSelectedGetColor)

	tagPathsPageMethodNormalGetStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodGetColor)

	tagPathsPageMethodSelectedPostStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodSelectedPostColor)

	tagPathsPageMethodNormalPostStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodPostColor)

	tagPathsPageMethodSelectedDeleteStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodSelectedDeleteColor)

	tagPathsPageMethodNormalDeleteStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodDeleteColor)

	tagPathsPageMethodSelectedPutStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodSelectedPutColor)

	tagPathsPageMethodNormalPutStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodPutColor)

	tagPathsPageMethodSelectedPatchStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodSelectedPatchColor)

	tagPathsPageMethodNormalPatchStyle = tagPathsPageMethodStyle.Copy().
						Foreground(httpMethodPatchColor)

	tagPathsPageMethodDeprecatedStyle = tagPathsPageMethodStyle.Copy().
						Strikethrough(true)

	tagPathsPageMethodSelectedDeprecatedStyle = tagPathsPageMethodDeprecatedStyle.Copy().
							Foreground(httpMethodSelectedDeprecatedColor)

	tagPathsPageMethodNormalDeprecatedStyle = tagPathsPageMethodDeprecatedStyle.Copy().
						Foreground(httpMethodDeprecatedColor)

	tagPathsPagePathDeprecatedStyle = lipgloss.NewStyle().
					Strikethrough(true)

	tagPathsPagePathSelectedDeprecatedStyle = tagPathsPagePathDeprecatedStyle.Copy().
						Foreground(selectedColor)

	tagPathsPagePathNormalDeprecatedStyle = tagPathsPagePathDeprecatedStyle.Copy().
						Foreground(lipgloss.Color("246"))
)

type tagPathsPageListItem struct {
	path *topi.Path
}

var _ list.Item = (*tagPathsPageListItem)(nil)

func (i tagPathsPageListItem) FilterValue() string {
	return i.path.UriPath
}

func (i tagPathsPageListItem) styledTitle(selected bool) string {
	var method, path string
	if i.path.Deprecated {
		if selected {
			method = tagPathsPageMethodSelectedDeprecatedStyle.Render(i.path.Method)
			path = tagPathsPagePathSelectedDeprecatedStyle.Render(i.path.UriPath)
		} else {
			method = tagPathsPageMethodNormalDeprecatedStyle.Render(i.path.Method)
			path = tagPathsPagePathNormalDeprecatedStyle.Render(i.path.UriPath)
		}
		return fmt.Sprintf("%s %s", padding.String(method, 7), path)
	}
	switch i.path.Method {
	case http.MethodGet:
		if selected {
			method = tagPathsPageMethodSelectedGetStyle.Render(i.path.Method)
		} else {
			method = tagPathsPageMethodNormalGetStyle.Render(i.path.Method)
		}
	case http.MethodPost:
		if selected {
			method = tagPathsPageMethodSelectedPostStyle.Render(i.path.Method)
		} else {
			method = tagPathsPageMethodNormalPostStyle.Render(i.path.Method)
		}
	case http.MethodDelete:
		if selected {
			method = tagPathsPageMethodSelectedDeleteStyle.Render(i.path.Method)
		} else {
			method = tagPathsPageMethodNormalDeleteStyle.Render(i.path.Method)
		}
	case http.MethodPut:
		if selected {
			method = tagPathsPageMethodSelectedPutStyle.Render(i.path.Method)
		} else {
			method = tagPathsPageMethodNormalPutStyle.Render(i.path.Method)
		}
	case http.MethodPatch:
		if selected {
			method = tagPathsPageMethodSelectedPatchStyle.Render(i.path.Method)
		} else {
			method = tagPathsPageMethodNormalPatchStyle.Render(i.path.Method)
		}
	default:
		method = tagPathsPageMethodStyle.Render(i.path.Method)
	}
	if selected {
		path = listSelectedTitleColorStyle.Render(i.path.UriPath)
	} else {
		path = listNormalTitleColorStyle.Render(i.path.UriPath)
	}
	return fmt.Sprintf("%s %s", padding.String(method, 7), path)
}

func (i tagPathsPageListItem) styledDesc(selected bool, width int) string {
	desc := truncateWithTail(i.path.Summary, uint(width))
	if selected {
		desc = listSelectedDescColorStyle.Render(desc)
	} else {
		desc = listNormalDescColorStyle.Render(desc)
	}
	return desc
}

type tagPathsPageListDelegate struct{}

var _ list.ItemDelegate = (*tagPathsPageListDelegate)(nil)

func newTagPathsPageListDelegate() tagPathsPageListDelegate {
	return tagPathsPageListDelegate{}
}

func (d tagPathsPageListDelegate) Height() int {
	return 2
}

func (d tagPathsPageListDelegate) Spacing() int {
	return 1
}

func (d tagPathsPageListDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd {
	return nil
}

func (d tagPathsPageListDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i := item.(tagPathsPageListItem)
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
