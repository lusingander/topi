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

func containsString(v string, ss []string) bool {
	for _, s := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func ptr[T any](v T) *T {
	return &v
}

func digit(n uint) uint {
	if n == 0 {
		return 1
	}
	var c uint
	for n > 0 {
		n /= 10
		c++
	}
	return c
}
