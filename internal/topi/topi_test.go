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
