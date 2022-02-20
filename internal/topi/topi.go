package topi

import (
	"net/http"
	"sort"
	"strings"
)

const (
	UntaggedDummyTag = "<<untagged>>"
)

type Document struct {
	Info       *Info
	TagPathMap map[string][]*Path
	Tags       []*Tag
}

func NewDocument(info *Info, tagPathMap map[string][]*Path, tags []*Tag) *Document {
	for _, paths := range tagPathMap {
		sortPaths(paths)
	}
	tags = mergeTags(tagPathMap, tags)
	sortTags(tags)
	return &Document{
		Info:       info,
		TagPathMap: tagPathMap,
		Tags:       tags,
	}
}

func sortPaths(paths []*Path) {
	sort.Slice(paths, func(i, j int) bool { return comparePath(paths[i], paths[j]) })
}

func mergeTags(tagPathMap map[string][]*Path, tags []*Tag) []*Tag {
	ret := tags
	for tagName := range tagPathMap {
		if !eixstTag(tagName, tags) {
			tag := &Tag{
				Name: tagName,
			}
			ret = append(ret, tag)
		}
	}
	return ret
}

func eixstTag(name string, tags []*Tag) bool {
	for _, tag := range tags {
		if name == tag.Name {
			return true
		}
	}
	return false
}

func sortTags(tags []*Tag) {
	sort.Slice(tags, func(i, j int) bool { return compareTags(tags[i].Name, tags[j].Name) })
}

func compareTags(t1, t2 string) bool {
	if t1 == UntaggedDummyTag {
		return false
	}
	if t2 == UntaggedDummyTag {
		return true
	}
	return t1 < t2
}

func (d *Document) FindPathByOperationId(operationId string) *Path {
	for _, paths := range d.TagPathMap {
		for _, path := range paths {
			if path.OperationId == operationId {
				return path
			}
		}
	}
	return nil
}

type Info struct {
	OpenAPIVersion    string
	Title             string
	Description       string
	TermsOfService    string
	ContactName       string
	ContactUrl        string
	ContactEmail      string
	LicenseName       string
	LicenseUrl        string
	Version           string
	ExDocsDescription string
	ExDocsUrl         string
}

type Path struct {
	UriPath          string
	Method           string
	OperationId      string
	Summary          string
	Description      string
	Deprecated       bool
	PathParameters   []*Parameter
	QueryParameters  []*Parameter
	HeaderParameters []*Parameter
	CookieParameters []*Parameter
}

func comparePath(p1, p2 *Path) bool {
	p1Paths := strings.Split(p1.UriPath, "/")
	p2Paths := strings.Split(p2.UriPath, "/")
	len1 := len(p1Paths)
	len2 := len(p2Paths)
	if len1 < len2 {
		return true
	}
	if len1 > len2 {
		return false
	}
	for i := 0; i < len1; i++ {
		if p1Paths[i] != p2Paths[i] {
			return comparePathElem(p1Paths[i], p2Paths[i])
		}
	}
	return httpMethod[p1.Method] < httpMethod[p2.Method]
}

func comparePathElem(e1, e2 string) bool {
	b1 := strings.HasPrefix(e1, "{")
	b2 := strings.HasPrefix(e2, "{")
	if b1 && !b2 {
		return true
	}
	if !b1 && b2 {
		return false
	}
	return e1 < e2
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

type Parameter struct {
	Name        string
	In          string
	Description string
	Required    bool
	Deprecated  bool
	Schema      *Schema
}

type Schema struct {
	Type   string
	Format string
}

type Tag struct {
	Name        string
	Description string
}
