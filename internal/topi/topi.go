package topi

import (
	"net/http"
	"sort"
	"strings"
)

type Document struct {
	TagPathMap map[string][]*Path
	Tags       []*Tag
}

func NewDocument(tagPathMap map[string][]*Path, tags []*Tag) *Document {
	for _, paths := range tagPathMap {
		sortPaths(paths)
	}
	sortTags(tags)
	return &Document{
		TagPathMap: tagPathMap,
		Tags:       tags,
	}
}

func sortPaths(paths []*Path) {
	sort.Slice(paths, func(i, j int) bool { return comparePath(paths[i], paths[j]) })
}

func sortTags(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool { return tags[i].Name < tags[j].Name })
}

type Path struct {
	UriPath     string
	Method      string
	OperationId string
	Summary     string
	Deprecated  bool
}

func comparePath(p1, p2 *Path) bool {
	// todo: fix logic
	p1Paths := strings.Split(p1.UriPath, "/")
	p2Paths := strings.Split(p2.UriPath, "/")
	if len(p1Paths) < len(p2Paths) {
		return true
	}
	if len(p1Paths) > len(p2Paths) {
		return false
	}
	for i := 0; i < len(p1Paths); i++ {
		if p1Paths[i] != p2Paths[i] {
			return p1Paths[i] < p2Paths[i]
		}
	}
	return httpMethod[p1.Method] < httpMethod[p2.Method]
}

var httpMethod = map[string]int{
	http.MethodGet:     0,
	http.MethodHead:    1,
	http.MethodPost:    2,
	http.MethodPut:     3,
	http.MethodPatch:   4,
	http.MethodDelete:  5,
	http.MethodConnect: 6,
	http.MethodOptions: 7,
	http.MethodTrace:   8,
}

type Tag struct {
	Name        string
	Description string
}
