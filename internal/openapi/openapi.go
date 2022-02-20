package openapi

import (
	"context"
	"net/url"
	"path/filepath"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/lusingander/topi/internal/topi"
)

func Load(path string) (*topi.Document, error) {
	ctx := context.Background()
	loader := openapi3.Loader{
		Context:               ctx,
		IsExternalRefsAllowed: true,
	}

	if uri, err := url.ParseRequestURI(path); err == nil {
		if doc, err := loader.LoadFromURI(uri); err == nil {
			return convert(doc), nil
		}
	}

	fp, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	doc, err := loader.LoadFromFile(fp)
	if err != nil {
		return nil, err
	}
	return convert(doc), nil
}

func convert(t *openapi3.T) *topi.Document {
	info := convertInfo(t.OpenAPI, t.Info, t.ExternalDocs)
	paths := convertPaths(t.Paths)
	tags := convertTags(t.Tags)
	return topi.NewDocument(info, paths, tags)
}

func convertInfo(openapi string, info *openapi3.Info, exDocs *openapi3.ExternalDocs) *topi.Info {
	contactName, contactUrl, contactEmail := convertInfoContact(info.Contact)
	licenseName, licenseUrl := convertInfoLicense(info.License)
	exDocsDesc, exDocsUrl := convertExternalDocs(exDocs)
	return &topi.Info{
		OpenAPIVersion:    openapi,
		Title:             info.Title,
		Description:       info.Description,
		TermsOfService:    info.TermsOfService,
		ContactName:       contactName,
		ContactUrl:        contactUrl,
		ContactEmail:      contactEmail,
		LicenseName:       licenseName,
		LicenseUrl:        licenseUrl,
		Version:           info.Version,
		ExDocsDescription: exDocsDesc,
		ExDocsUrl:         exDocsUrl,
	}
}

func convertInfoContact(c *openapi3.Contact) (name, url, email string) {
	if c == nil {
		name, url, email = "", "", ""
	} else {
		name, url, email = c.Name, c.URL, c.Email
	}
	return
}

func convertInfoLicense(l *openapi3.License) (name, url string) {
	if l == nil {
		name, url = "", ""
	} else {
		name, url = l.Name, l.URL
	}
	return
}

func convertExternalDocs(d *openapi3.ExternalDocs) (desc, url string) {
	if d == nil {
		desc, url = "", ""
	} else {
		desc, url = d.Description, d.URL
	}
	return
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
		path := convertOperation(pathItem, op, method, uriPath)
		tag := getTag(op)
		if _, ok := ret[tag]; !ok {
			ret[tag] = make([]*topi.Path, 0)
		}
		ret[tag] = append(ret[tag], path)
	}
	return ret
}

func getTag(op *openapi3.Operation) string {
	if len(op.Tags) == 0 {
		return topi.UntaggedDummyTag
	}
	return op.Tags[0]
}

func convertOperation(pathItem *openapi3.PathItem, op *openapi3.Operation, method, uriPath string) *topi.Path {
	params := op.Parameters
	if len(params) == 0 {
		params = pathItem.Parameters
	}

	ret := &topi.Path{
		UriPath:          uriPath,
		Method:           method,
		OperationId:      op.OperationID,
		Summary:          op.Summary,
		Description:      op.Description,
		Deprecated:       op.Deprecated,
		PathParameters:   convertParameters(params, "path"),
		QueryParameters:  convertParameters(params, "query"),
		HeaderParameters: convertParameters(params, "header"),
		CookieParameters: convertParameters(params, "cookie"),
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

func convertParameters(params openapi3.Parameters, in string) []*topi.Parameter {
	ret := make([]*topi.Parameter, 0)
	for _, param := range params {
		if param.Value.In == in {
			p := &topi.Parameter{
				Name:        param.Value.Name,
				In:          param.Value.In,
				Description: param.Value.Description,
				Required:    param.Value.Required,
				Deprecated:  param.Value.Deprecated,
				Schema:      convertSchema(param.Value.Schema),
			}
			ret = append(ret, p)
		}
	}
	return ret
}

func convertSchema(schema *openapi3.SchemaRef) *topi.Schema {
	if schema == nil || schema.Value == nil {
		return nil
	}
	sc := schema.Value
	return &topi.Schema{
		Type:         sc.Type,
		Format:       sc.Format,
		Default:      sc.Default,
		Enum:         sc.Enum,
		Min:          sc.Min,
		Max:          sc.Max,
		ExclusiveMin: sc.ExclusiveMin,
		ExclusiveMax: sc.ExclusiveMax,
		MultipleOf:   sc.MultipleOf,
		MinLength:    sc.MinLength,
		MaxLength:    sc.MaxLength,
		Pattern:      sc.Pattern,
		MinItems:     sc.MinItems,
		MaxItems:     sc.MaxItems,
		Items:        convertSchema(sc.Items),
	}
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
	return ret
}
