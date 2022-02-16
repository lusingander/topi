package ui

import (
	"strings"

	"github.com/muesli/reflow/truncate"
)

func truncateWithTail(s string, w uint) string {
	tail := "..."
	ss := strings.ReplaceAll(s, "\n", " ")
	return truncate.StringWithTail(ss, w, tail)
}
