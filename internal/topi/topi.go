package topi

import (
	"net/http"
	"sort"
	"strings"
)

const (
	UntaggedDummyTag = "<untagged>"
)

type Document struct {
	Meta       *Meta
	Info       *Info
	TagPathMap map[string][]*Path
	Tags       []*Tag
	Components *Components
}

func NewDocument(meta *Meta, info *Info, tagPathMap map[string][]*Path, tags []*Tag, components *Components) *Document {
	for _, paths := range tagPathMap {
		sortPaths(paths)
	}
	tags = mergeTags(tagPathMap, tags)
	sortTags(tags)
	return &Document{
		Meta:       meta,
		Info:       info,
		TagPathMap: tagPathMap,
		Tags:       tags,
		Components: components,
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

type Meta struct {
	FileName string
	FullPath string
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
	RequestBody      *RequestBody
	Responses        []*Response
	Security         []*SecurityRequirement
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
	Type        string
	Format      string
	Default     interface{}
	Enum        []interface{}
	Description string
	Deprecated  bool
	ReadOnly    bool
	WriteOnly   bool
	OneOf       []*Schema
	AllOf       []*Schema

	// number/integer
	Min          *float64
	Max          *float64
	ExclusiveMin bool
	ExclusiveMax bool
	MultipleOf   *float64

	// string
	MinLength uint64
	MaxLength *uint64
	Pattern   string

	// array
	MinItems uint64
	MaxItems *uint64
	Items    *Schema

	// object
	Required   []string
	Properties map[string]*Schema
}

func (s *Schema) MergedAllOf() *Schema {
	if len(s.AllOf) == 0 {
		return nil
	}
	ret := &Schema{}
	for _, schema := range s.AllOf {
		if schema.Type != "" {
			ret.Type = schema.Type
		}
		if len(schema.Properties) > 0 {
			if ret.Properties == nil {
				ret.Properties = make(map[string]*Schema)
			}
			for k, v := range schema.Properties {
				if prop, ok := ret.Properties[k]; ok {
					sc := &Schema{AllOf: []*Schema{prop, v}}
					ret.Properties[k] = sc.MergedAllOf()
				} else {
					ret.Properties[k] = v
				}
			}
		}
		if len(schema.Required) > 0 {
			ret.Required = append(ret.Required, schema.Required...)
		}
	}
	return ret
}

type RequestBody struct {
	Description string
	Required    bool
	Conetnt     []*MediaTypeContent
}

type MediaTypeContent struct {
	MediaType string
	Schema    *Schema
}

type Response struct {
	StatusCode  string
	Description string
	Conetnt     []*MediaTypeContent
	Headers     []*Header
}

type Header struct {
	Name      string
	Parameter *Parameter
}

func (r *Response) Success() bool {
	return strings.HasPrefix(r.StatusCode, "2")
}

func (r *Response) Error() bool {
	return strings.HasPrefix(r.StatusCode, "4") || strings.HasPrefix(r.StatusCode, "5")
}

type Tag struct {
	Name        string
	Description string
}

type SecurityRequirement struct {
	Schemes []*SecurityRequirementScheme
}

type SecurityRequirementScheme struct {
	Key    string
	Scopes []string
}

type Components struct {
	SecuritySchemes []*SecurityScheme
}

type SecurityScheme struct {
	Key         string
	Type        string
	Description string

	Name string // api_key
	In   string // api_key

	Scheme       string // http
	BearerFormat string // http

	OpenIdConnectUrl string // openIdConnect

	OAuthFlows *OAtuhFlows // oauth2
}

func (s *SecurityScheme) TypeStr() string {
	switch s.Type {
	case "apiKey":
		return "API Key"
	case "http":
		return "HTTP"
	case "oauth2":
		return "OAuth2"
	case "openIdConnect":
		return "OpenID Connect"
	default:
		return ""
	}
}

type OAtuhFlows struct {
	ImplicitFlow                         *OAuthFlow
	ResourceOwnerPasswordCredentialsFlow *OAuthFlow
	ClientCredentialsFlow                *OAuthFlow
	AuthorizatonCodeFlow                 *OAuthFlow
}

type OAuthFlow struct {
	AuthorizationURL string
	TokenURL         string
	RefreshURL       string
	Scopes           []*Scope
}

type Scope struct {
	Name   string
	Detail string
}
