package ui

import (
	"fmt"
	"strings"

	"github.com/lusingander/topi/internal/topi"
)

func schemaTypeString(sc *topi.Schema) string {
	var s strings.Builder
	if sc.Type != "" {
		if sc.Type == "array" {
			itemType := sc.Items.Type // schema.items must be present if the type is array
			s.WriteString(fmt.Sprintf("array of %s", itemType))
		} else {
			s.WriteString(sc.Type)
			if sc.Format != "" {
				s.WriteString(fmt.Sprintf("(%s)", sc.Format))
			}
		}
	}
	return s.String()
}

func schemaConstraintStrings(sc *topi.Schema) []string {
	ret := make([]string, 0)
	switch sc.Type {
	case "integer", "number":
		if sc.Min != nil || sc.Max != nil {
			s := "n"
			if sc.Min != nil {
				if sc.ExclusiveMin {
					s = fmt.Sprintf("%g < %s", *sc.Min, s)
				} else {
					s = fmt.Sprintf("%g <= %s", *sc.Min, s)
				}
			}
			if sc.Max != nil {
				if sc.ExclusiveMax {
					s = fmt.Sprintf("%s < %g", s, *sc.Max)
				} else {
					s = fmt.Sprintf("%s <= %g", s, *sc.Max)
				}
			}
			ret = append(ret, s)
		}
		if sc.MultipleOf != nil {
			s := fmt.Sprintf("multiple of %g", *sc.MultipleOf)
			ret = append(ret, s)
		}
	case "string":
		if sc.MinLength > 0 || sc.MaxLength != nil {
			s := "len"
			if sc.MinLength > 0 {
				s = fmt.Sprintf("%d <= %s", sc.MinLength, s)
			}
			if sc.MaxLength != nil {
				s = fmt.Sprintf("%s <= %d", s, *sc.MaxLength)
			}
			ret = append(ret, s)
		}
		if sc.Pattern != "" {
			s := sc.Pattern
			ret = append(ret, s)
		}
	case "array":
		if sc.MinItems > 0 || sc.MaxItems != nil {
			s := "items"
			if sc.MinItems > 0 {
				s = fmt.Sprintf("%d <= %s", sc.MinItems, s)
			}
			if sc.MaxItems != nil {
				s = fmt.Sprintf("%s <= %d", s, *sc.MaxItems)
			}
			ret = append(ret, s)
		}
	case "boolean", "object":
		// do nothing
	}
	return ret
}

func sliceString(vs []interface{}) string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = fmt.Sprintf("%v", v)
	}
	return fmt.Sprintf("[%s]", strings.Join(ss, ", "))
}
