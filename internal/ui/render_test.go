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
	}

	for _, test := range tests {
		got := schemaConstraintStrings(test.schema)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got=%v, want=%v", got, test.want)
		}
	}
}
