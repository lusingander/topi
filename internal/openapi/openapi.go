package openapi

import (
	"context"
	"sort"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/lusingander/topi/internal/topi"
)

func Load(filepath string) (*topi.Document, error) {
	ctx := context.Background()
	loader := openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}
	doc, err := loader.LoadFromFile(filepath)
	if err != nil {
		return nil, err
	}
	return convert(doc), nil
}

func convert(t *openapi3.T) *topi.Document {
	ret := &topi.Document{
		TagPathMap: convertPaths(t.Paths),
		Tags:       convertTags(t.Tags),
	}
	return ret
}

func convertPaths(paths openapi3.Paths) map[string][]*topi.Path {
	ret := make(map[string][]*topi.Path)
	for k, v := range paths {
		items := convertPathItem(v, k)
		ret = mergeMap(ret, items)
	}
	return ret
}

func convertPathItem(pathItem *openapi3.PathItem, uriPath string) map[string][]*topi.Path {
	ret := make(map[string][]*topi.Path)
	for method, op := range pathItem.Operations() {
		path := convertOperation(op, method, uriPath)
		tag := getTag(op)
		if _, ok := ret[tag]; !ok {
			ret[tag] = make([]*topi.Path, 0)
		}
		ret[tag] = append(ret[tag], path)
	}
	return ret
}

func getTag(op *openapi3.Operation) string {
	return op.Tags[0]
}

func convertOperation(op *openapi3.Operation, method, uriPath string) *topi.Path {
	ret := &topi.Path{
		UriPath:     uriPath,
		Method:      method,
		OperationId: op.OperationID,
		Summary:     op.Summary,
		Deprecated:  op.Deprecated,
	}
	return ret
}

func mergeMap(m1, m2 map[string][]*topi.Path) map[string][]*topi.Path {
	ret := make(map[string][]*topi.Path)
	for k, v := range m1 {
		if _, ok := ret[k]; ok {
			ret[k] = append(ret[k], v...)
		} else {
			ret[k] = v
		}
	}
	for k, v := range m2 {
		if _, ok := ret[k]; ok {
			ret[k] = append(ret[k], v...)
		} else {
			ret[k] = v
		}
	}
	return ret
}

func convertTags(tags openapi3.Tags) []*topi.Tag {
	ret := make([]*topi.Tag, 0)
	for _, tag := range tags {
		t := &topi.Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}
		ret = append(ret, t)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i].Name < ret[j].Name })
	return ret
}
