package openapi

import (
	"reflect"
	"testing"

	"github.com/lusingander/topi/internal/topi"
)

func TestMergeMap(t *testing.T) {
	m1 := map[string][]*topi.Path{
		"foo": {
			&topi.Path{UriPath: "/abc", Method: "GET"},
			&topi.Path{UriPath: "/abc", Method: "POST"},
		},
		"bar": {
			&topi.Path{UriPath: "/stu", Method: "GET"},
			&topi.Path{UriPath: "/stu/{id}", Method: "GET"},
		},
	}
	m2 := map[string][]*topi.Path{
		"foo": {
			&topi.Path{UriPath: "/abc", Method: "DELETE"},
			&topi.Path{UriPath: "/abc/{id}", Method: "GET"},
		},
		"baz": {
			&topi.Path{UriPath: "/xyz", Method: "POST"},
		},
	}
	want := map[string][]*topi.Path{
		"foo": {
			&topi.Path{UriPath: "/abc", Method: "GET"},
			&topi.Path{UriPath: "/abc", Method: "POST"},
			&topi.Path{UriPath: "/abc", Method: "DELETE"},
			&topi.Path{UriPath: "/abc/{id}", Method: "GET"},
		},
		"bar": {
			&topi.Path{UriPath: "/stu", Method: "GET"},
			&topi.Path{UriPath: "/stu/{id}", Method: "GET"},
		},
		"baz": {
			&topi.Path{UriPath: "/xyz", Method: "POST"},
		},
	}
	got := mergeMap(m1, m2)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got=%v, want=%v", got, want)
	}
}
