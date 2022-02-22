package ui

import (
	"reflect"
	"testing"

	"github.com/lusingander/topi/internal/topi"
)

func TestSchemaConstraintStrings(t *testing.T) {
	tests := []struct {
		schema *topi.Schema
		want   []string
	}{
		{
			schema: &topi.Schema{
				Type: "number",
				Min:  float64Pointer(1),
				Max:  float64Pointer(10),
			},
			want: []string{
				"1 <= n <= 10",
			},
		},
		{
			schema: &topi.Schema{
				Type:         "number",
				Min:          float64Pointer(-123.45),
				Max:          float64Pointer(0),
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
				MultipleOf: float64Pointer(5),
			},
			want: []string{
				"multiple of 5",
			},
		},
		{
			schema: &topi.Schema{
				Type:       "number",
				Max:        float64Pointer(20),
				MultipleOf: float64Pointer(0.5),
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
				MaxLength: uint64Pointer(30),
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
				MaxItems: uint64Pointer(5),
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