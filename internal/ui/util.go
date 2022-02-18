package ui

import (
	"strings"

	"github.com/muesli/reflow/truncate"
	"github.com/pkg/browser"
)

func truncateWithTail(s string, w uint) string {
	tail := "..."
	ss := strings.ReplaceAll(s, "\n", " ")
	return truncate.StringWithTail(ss, w, tail)
}

func openInBrowser(url string) error {
	return browser.OpenURL(url)
}
