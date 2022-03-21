package topi

import (
	"reflect"
	"testing"
)

func TestSortPaths(t *testing.T) {
	paths := []*Path{
		{UriPath: "/abc", Method: "POST"},
		{UriPath: "/abc/{xId}/stu", Method: "POST"},
		{UriPath: "/abc/stu", Method: "DELETE"},
		{UriPath: "/abc", Method: "DELETE"},
		{UriPath: "/def/{yId}", Method: "PUT"},
		{UriPath: "/def", Method: "POST"},
		{UriPath: "/abc/{xId}/stu/{yId}", Method: "GET"},
		{UriPath: "/abc", Method: "PUT"},
		{UriPath: "/abc/{xId}", Method: "GET"},
		{UriPath: "/abc", Method: "GET"},
		{UriPath: "/def/{xId}", Method: "POST"},
		{UriPath: "/abc", Method: "PATCH"},
		{UriPath: "/abc/xyz/{xId}", Method: "PATCH"},
		{UriPath: "/def", Method: "GET"},
	}
	want := []*Path{
		{UriPath: "/abc", Method: "GET"},
		{UriPath: "/abc", Method: "POST"},
		{UriPath: "/abc", Method: "PUT"},
		{UriPath: "/abc", Method: "PATCH"},
		{UriPath: "/abc", Method: "DELETE"},
		{UriPath: "/def", Method: "GET"},
		{UriPath: "/def", Method: "POST"},
		{UriPath: "/abc/{xId}", Method: "GET"},
		{UriPath: "/abc/stu", Method: "DELETE"},
		{UriPath: "/def/{xId}", Method: "POST"},
		{UriPath: "/def/{yId}", Method: "PUT"},
		{UriPath: "/abc/{xId}/stu", Method: "POST"},
		{UriPath: "/abc/xyz/{xId}", Method: "PATCH"},
		{UriPath: "/abc/{xId}/stu/{yId}", Method: "GET"},
	}
	sortPaths(paths)
	if !reflect.DeepEqual(paths, want) {
		t.Errorf("got=%v, want=%v", paths, want)
	}
}

func TestMergeTags(t *testing.T) {
	tagPathMap := map[string][]*Path{
		"foo": {
			{UriPath: "/abc", Method: "GET"},
			{UriPath: "/abc", Method: "POST"},
		},
		"bar": {
			{UriPath: "/def/{xId}", Method: "POST"},
		},
	}
	tags := []*Tag{
		{"foo", "foo detail"},
		{"baz", "baz detail"},
	}
	want := []*Tag{
		{"foo", "foo detail"},
		{"baz", "baz detail"},
		{"bar", ""},
	}
	got := mergeTags(tagPathMap, tags)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got=%v, want=%v", got, want)
	}
}

func TestMergedAllOf(t *testing.T) {
	tests := []struct {
		schema *Schema
		want   *Schema
	}{
		{
			schema: &Schema{
				AllOf: []*Schema{
					{
						Type: "integer",
					},
					{
						Type: "string",
					},
				},
			},
			want: &Schema{
				Type: "string",
			},
		},
		{
			schema: &Schema{
				AllOf: []*Schema{
					{
						Type: "integer",
					},
					{
						Type: "object",
						Properties: map[string]*Schema{
							"foo": {
								Type: "string",
							},
						},
					},
				},
			},
			want: &Schema{
				Type: "object",
				Properties: map[string]*Schema{
					"foo": {
						Type: "string",
					},
				},
			},
		},
		{
			schema: &Schema{
				AllOf: []*Schema{
					{
						Type: "object",
						Properties: map[string]*Schema{
							"foo": {
								Type:   "integer",
								Format: "int32",
							},
							"bar": {
								Type: "integer",
							},
						},
						Required: []string{"foo"},
					},
					{
						Type: "object",
						Properties: map[string]*Schema{
							"bar": {
								Type: "string",
							},
							"baz": {
								Type:   "string",
								Format: "password",
							},
						},
						Required: []string{"bar"},
					},
				},
			},
			want: &Schema{
				Type: "object",
				Properties: map[string]*Schema{
					"foo": {
						Type:   "integer",
						Format: "int32",
					},
					"bar": {
						Type: "string",
					},
					"baz": {
						Type:   "string",
						Format: "password",
					},
				},
				Required: []string{"foo", "bar"},
			},
		},
		{
			schema: &Schema{
				AllOf: []*Schema{
					{
						Type: "object",
						Properties: map[string]*Schema{
							"foo": {
								Type:         "integer",
								Min:          ptr[float64](10),
								Max:          ptr[float64](30),
								ExclusiveMin: true,
								ExclusiveMax: true,
								MultipleOf:   ptr[float64](5),
							},
						},
					},
					{
						Type: "object",
						Properties: map[string]*Schema{
							"bar": {
								Type:      "string",
								MinLength: 8,
								MaxLength: ptr[uint64](32),
								Pattern:   "/(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])/",
							},
						},
					},
				},
			},
			want: &Schema{
				Type: "object",
				Properties: map[string]*Schema{
					"foo": {
						Type:         "integer",
						Min:          ptr[float64](10),
						Max:          ptr[float64](30),
						ExclusiveMin: true,
						ExclusiveMax: true,
						MultipleOf:   ptr[float64](5),
					},
					"bar": {
						Type:      "string",
						MinLength: 8,
						MaxLength: ptr[uint64](32),
						Pattern:   "/(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])/",
					},
				},
			},
		},
	}

	for _, test := range tests {
		got := test.schema.MergedAllOf()
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("got=%v, want=%v", got, test.want)
		}
	}
}

func TestMergedAllOf_NoSideEffect(t *testing.T) {
	schema := &Schema{
		AllOf: []*Schema{
			{
				Type:       "object",
				Properties: map[string]*Schema{"foo": {Type: "integer"}},
			},
			{
				Type:       "object",
				Properties: map[string]*Schema{"bar": {Type: "string"}},
			},
		},
	}
	want := &Schema{
		Type:       "object",
		Properties: map[string]*Schema{"foo": {Type: "integer"}, "bar": {Type: "string"}},
	}
	schemaWant := &Schema{
		AllOf: []*Schema{
			{
				Type:       "object",
				Properties: map[string]*Schema{"foo": {Type: "integer"}},
			},
			{
				Type:       "object",
				Properties: map[string]*Schema{"bar": {Type: "string"}},
			},
		},
	}
	got := schema.MergedAllOf()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got=%v, want=%v", got, want)
	}
	if !reflect.DeepEqual(schema, schemaWant) {
		t.Errorf("before=%v, after=%v", schema, schemaWant)
	}
}

func ptr[T any](v T) *T {
	return &v
}
