package ui

import (
	"reflect"
	"testing"

	"github.com/lusingander/topi/internal/topi"
)

func TestSchemaTypeString(t *testing.T) {
	tests := []struct {
		schema *topi.Schema
		want   string
	}{
		{
			schema: &topi.Schema{},
			want:   "",
		},
		{
			schema: &topi.Schema{
				Type: "integer",
			},
			want: "integer",
		},
		{
			schema: &topi.Schema{
				Type:   "string",
				Format: "password",
			},
			want: "string(password)",
		},
		{
			schema: &topi.Schema{
				Type: "array",
				Items: &topi.Schema{
					Type: "boolean",
				},
			},
			want: "array of boolean",
		},
		{
			schema: &topi.Schema{
				Type: "object",
				Properties: map[string]*topi.Schema{
					"foo": {Type: "number"},
				},
			},
			want: "object",
		},
		{
			schema: &topi.Schema{
				OneOf: []*topi.Schema{
					{Type: "integer"},
					{Type: "object", Properties: map[string]*topi.Schema{}},
					{Type: "string"},
					{Type: "object", Properties: map[string]*topi.Schema{}},
				},
			},
			want: "one of (integer | object[2] | string | object[4])",
		},
		{
			schema: &topi.Schema{
				AllOf: []*topi.Schema{
					{Type: "integer"},
					{Type: "boolean"},
					{Type: "string"},
				},
			},
			want: "string",
		},
		{
			schema: &topi.Schema{
				AllOf: []*topi.Schema{
					{Type: "object", Properties: map[string]*topi.Schema{}},
					{Type: "object", Properties: map[string]*topi.Schema{}},
				},
			},
			want: "object",
		},
	}
	for _, test := range tests {
		got := schemaTypeString(test.schema)
		if got != test.want {
			t.Errorf("got=%v, want=%v", got, test.want)
		}
	}
}

func TestSchemaConstraintStrings(t *testing.T) {
	tests := []struct {
		schema *topi.Schema
		want   []string
	}{
		{
			schema: &topi.Schema{
				Type: "number",
				Min:  ptr[float64](1),
				Max:  ptr[float64](10),
			},
			want: []string{
				"1 <= n <= 10",
			},
		},
		{
			schema: &topi.Schema{
				Type:         "number",
				Min:          ptr(-123.45),
				Max:          ptr[float64](0),
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			want: []string{
				"-123.45 < n < 0",
			},
		},
		{
			schema: &topi.Schema{
				Type:       "integer",
				MultipleOf: ptr[float64](5),
			},
			want: []string{
				"multiple of 5",
			},
		},
		{
			schema: &topi.Schema{
				Type:       "number",
				Max:        ptr[float64](20),
				MultipleOf: ptr(0.5),
			},
			want: []string{
				"n <= 20",
				"multiple of 0.5",
			},
		},
		{
			schema: &topi.Schema{
				Type:      "string",
				MinLength: 1,
				MaxLength: ptr[uint64](30),
				Pattern:   "/(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])/",
			},
			want: []string{
				"1 <= len <= 30",
				"/(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])/",
			},
		},
		{
			schema: &topi.Schema{
				Type:     "array",
				MinItems: 2,
				MaxItems: ptr[uint64](5),
			},
			want: []string{
				"2 <= items <= 5",
			},
		},
	}

	for _, test := range tests {
		got := schemaConstraintStrings(test.schema)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got=%v, want=%v", got, test.want)
		}
	}
}

func TestSliceString(t *testing.T) {
	tests := []struct {
		slice []interface{}
		want  string
	}{
		{
			slice: []interface{}{},
			want:  "[]",
		},
		{
			slice: []interface{}{"foo", "bar", "baz"},
			want:  "[foo, bar, baz]",
		},
		{
			slice: []interface{}{1},
			want:  "[1]",
		},
		{
			slice: []interface{}{0.25, "abc"},
			want:  "[0.25, abc]",
		},
	}

	for _, test := range tests {
		got := sliceString(test.slice)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got=%v, want=%v", got, test.want)
		}
	}
}
