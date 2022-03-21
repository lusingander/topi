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
			Margin:      ptr[uint](0),
			IndentToken: ptr(""),
			Indent:      ptr[uint](0),
		},
		Paragraph: ansi.StyleBlock{
			Margin:      ptr[uint](0),
			IndentToken: ptr(""),
			Indent:      ptr[uint](0),
			StylePrimitive: ansi.StylePrimitive{
				BlockPrefix: "",
				BlockSuffix: "",
				Prefix:      "",
			},
		},
		BlockQuote: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Color: ptr("246"),
			},
			Indent:      ptr[uint](1),
			IndentToken: ptr("| "),
		},
		List: ansi.StyleList{
			LevelIndent: 2,
		},
		Heading: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				BlockSuffix: "\n",
				Color:       ptr("70"),
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
			CrossedOut: ptr(true),
		},
		Emph: ansi.StylePrimitive{
			Italic: ptr(true),
		},
		Strong: ansi.StylePrimitive{
			Bold: ptr(true),
		},
		HorizontalRule: ansi.StylePrimitive{
			Color:  ptr("242"),
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
			Color:       ptr("33"),
			Underline:   ptr(true),
			BlockPrefix: "(",
			BlockSuffix: ")",
		},
		LinkText: ansi.StylePrimitive{
			Color: ptr("195"),
		},
		Image: ansi.StylePrimitive{
			Underline:   ptr(true),
			Color:       ptr("33"),
			BlockPrefix: "(",
			BlockSuffix: ")",
		},
		ImageText: ansi.StylePrimitive{
			Color: ptr("195"),
		},
		Code: ansi.StyleBlock{
			StylePrimitive: ansi.StylePrimitive{
				Prefix:          " ",
				Suffix:          " ",
				Color:           ptr("203"),
				BackgroundColor: ptr("236"),
			},
		},
		CodeBlock: ansi.StyleCodeBlock{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{
					Color: ptr("244"),
				},
				Margin: ptr[uint](0),
			},
			Chroma: &ansi.Chroma{
				Text: ansi.StylePrimitive{
					Color: ptr("#C4C4C4"),
				},
				Error: ansi.StylePrimitive{
					Color:           ptr("#F1F1F1"),
					BackgroundColor: ptr("#F05B5B"),
				},
				Comment: ansi.StylePrimitive{
					Color: ptr("#676767"),
				},
				CommentPreproc: ansi.StylePrimitive{
					Color: ptr("#FF875F"),
				},
				Keyword: ansi.StylePrimitive{
					Color: ptr("#00AAFF"),
				},
				KeywordReserved: ansi.StylePrimitive{
					Color: ptr("#FF5FD2"),
				},
				KeywordNamespace: ansi.StylePrimitive{
					Color: ptr("#FF5F87"),
				},
				KeywordType: ansi.StylePrimitive{
					Color: ptr("#6E6ED8"),
				},
				Operator: ansi.StylePrimitive{
					Color: ptr("#EF8080"),
				},
				Punctuation: ansi.StylePrimitive{
					Color: ptr("#E8E8A8"),
				},
				Name: ansi.StylePrimitive{
					Color: ptr("#C4C4C4"),
				},
				NameBuiltin: ansi.StylePrimitive{
					Color: ptr("#FF8EC7"),
				},
				NameTag: ansi.StylePrimitive{
					Color: ptr("#B083EA"),
				},
				NameAttribute: ansi.StylePrimitive{
					Color: ptr("#7A7AE6"),
				},
				NameClass: ansi.StylePrimitive{
					Color:     ptr("#F1F1F1"),
					Underline: ptr(true),
					Bold:      ptr(true),
				},
				NameDecorator: ansi.StylePrimitive{
					Color: ptr("#FFFF87"),
				},
				NameFunction: ansi.StylePrimitive{
					Color: ptr("#00D787"),
				},
				LiteralNumber: ansi.StylePrimitive{
					Color: ptr("#6EEFC0"),
				},
				LiteralString: ansi.StylePrimitive{
					Color: ptr("#C69669"),
				},
				LiteralStringEscape: ansi.StylePrimitive{
					Color: ptr("#AFFFD7"),
				},
				GenericDeleted: ansi.StylePrimitive{
					Color: ptr("#FD5B5B"),
				},
				GenericEmph: ansi.StylePrimitive{
					Italic: ptr(true),
				},
				GenericInserted: ansi.StylePrimitive{
					Color: ptr("#00D787"),
				},
				GenericStrong: ansi.StylePrimitive{
					Bold: ptr(true),
				},
				GenericSubheading: ansi.StylePrimitive{
					Color: ptr("#777777"),
				},
				Background: ansi.StylePrimitive{
					BackgroundColor: ptr("#373737"),
				},
			},
		},
		Table: ansi.StyleTable{
			StyleBlock: ansi.StyleBlock{
				StylePrimitive: ansi.StylePrimitive{},
			},
			CenterSeparator: ptr("+"),
			ColumnSeparator: ptr("|"),
			RowSeparator:    ptr("-"),
		},
		DefinitionDescription: ansi.StylePrimitive{
			BlockPrefix: " ",
		},
	}
)
