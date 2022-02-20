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

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}

func uintPtr(u uint) *uint {
	return &u
}

func float64Pointer(v float64) *float64 {
	return &v
}
