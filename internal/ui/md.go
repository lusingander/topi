package ui

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
)

func markdownRenderer(w int) (*glamour.TermRenderer, error) {
	return glamour.NewTermRenderer(
		glamour.WithStyles(glamourStyleConfig),
		glamour.WithWordWrap(w),
	)
}

var (
	glamourStyleConfig = ansi.StyleConfig{
		Document: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "",
				BlockSuffix: "",
				Prefix:      "",
			},
			Margin:      uintPtr(0),
			IndentToken: stringPtr(""),
			Indent:      uintPtr(0),
		},
		Paragraph: ansi.StyleBlock{
			Margin:      uintPtr(0),
			IndentToken: stringPtr(""),
			Indent:      uintPtr(0),
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "",
				BlockSuffix: "",
				Prefix:      "",
			},
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: stringPtr("246"),
			},
			Indent:      uintPtr(1),
			IndentToken: stringPtr("| "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "",
				Color:       stringPtr("70"),
			},
		},
		H1: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "# ",
			},
		},
		H2: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "## ",
			},
		},
		H3: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "### ",
			},
		},
		H4: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "#### ",
			},
		},
		H5: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "##### ",
			},
		},
		H6: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix: "###### ",
			},
		},
		Strikethrough: ansi.StylePrimitive{
			CrossedOut: boolPtr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: boolPtr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: boolPtr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  stringPtr("242"),
			Format: "----",
		},
		Item: ansi.StylePrimitive{
			BlockPrefix: "• ",
		},
		Enumeration: ansi.StylePrimitive{
			BlockPrefix: ". ",
		},
		Task: ansi.StyleTask{
			StylePrimitive: ansi.StylePrimitive{},
			Ticked:         "[x] ",
			Unticked:       "[ ] ",
		},
		Link: ansi.StylePrimitive{
			Color:       stringPtr("33"),
			Underline:   boolPtr(true),
			BlockPrefix: "(",
			BlockSuffix: ")",
		},
		LinkText: ansi.StylePrimitive{
			Color: stringPtr("195"),
		},
		Image: ansi.StylePrimitive{
			Underline:   boolPtr(true),
			Color:       stringPtr("33"),
			BlockPrefix: "(",
			BlockSuffix: ")",
		},
		ImageText: ansi.StylePrimitive{
			Color: stringPtr("195"),
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           stringPtr("203"),
				BackgroundColor: stringPtr("236"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: stringPtr("244"),
				},
				Margin: uintPtr(0),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: stringPtr("#C4C4C4"),
				},
				Error: ansi.StylePrimitive{
					Color:           stringPtr("#F1F1F1"),
					BackgroundColor: stringPtr("#F05B5B"),
				},
				Comment: ansi.StylePrimitive{
					Color: stringPtr("#676767"),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: stringPtr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					Color: stringPtr("#00AAFF"),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: stringPtr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: stringPtr("#FF5F87"),
				},
				KeywordType: ansi.StylePrimitive{
					Color: stringPtr("#6E6ED8"),
				},
				Operator: ansi.StylePrimitive{
					Color: stringPtr("#EF8080"),
				},
				Punctuation: ansi.StylePrimitive{
					Color: stringPtr("#E8E8A8"),
				},
				Name: ansi.StylePrimitive{
					Color: stringPtr("#C4C4C4"),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: stringPtr("#FF8EC7"),
				},
				NameTag: ansi.StylePrimitive{
					Color: stringPtr("#B083EA"),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: stringPtr("#7A7AE6"),
				},
				NameClass: ansi.StylePrimitive{
					Color:     stringPtr("#F1F1F1"),
					Underline: boolPtr(true),
					Bold:      boolPtr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: stringPtr("#FFFF87"),
				},
				NameFunction: ansi.StylePrimitive{
					Color: stringPtr("#00D787"),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: stringPtr("#6EEFC0"),
				},
				LiteralString: ansi.StylePrimitive{
					Color: stringPtr("#C69669"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: stringPtr("#AFFFD7"),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: stringPtr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: boolPtr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: stringPtr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: boolPtr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: stringPtr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: stringPtr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: stringPtr("+"),
			ColumnSeparator: stringPtr("|"),
			RowSeparator:    stringPtr("-"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: " ",
		},
	}
)
