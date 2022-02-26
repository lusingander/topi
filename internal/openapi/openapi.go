package openapi

import (
	"context"
	"net/url"
	"path"
	"path/filepath"
	"sort"

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
			return convert(path, doc), nil
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
	return convert(fp, doc), nil
}

func convert(filepath string, t *openapi3.T) *topi.Document {
	meta := convertMeta(filepath)
	info := convertInfo(t.OpenAPI, t.Info, t.ExternalDocs)
	paths := convertPaths(t.Paths)
	tags := convertTags(t.Tags)
	components := convertComponents(&t.Components)
	return topi.NewDocument(meta, info, paths, tags, components)
}

func convertMeta(filepath string) *topi.Meta {
	return &topi.Meta{
		FileName: path.Base(filepath),
		FullPath: filepath,
	}
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
		RequestBody:      convertRequestBody(op.RequestBody),
		Responses:        convertResponses(op.Responses),
		Security:         convertSecurityRequirements(op.Security),
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
		Description:  sc.Description,
		Deprecated:   sc.Deprecated,
		ReadOnly:     sc.ReadOnly,
		WriteOnly:    sc.WriteOnly,
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
		Required:     sc.Required,
		Properties:   convertSchemas(sc.Properties),
	}
}

func convertSchemas(s openapi3.Schemas) map[string]*topi.Schema {
	ret := make(map[string]*topi.Schema)
	for k, v := range s {
		ret[k] = convertSchema(v)
	}
	return ret
}

func convertRequestBody(body *openapi3.RequestBodyRef) *topi.RequestBody {
	if body == nil || body.Value == nil {
		return nil
	}
	b := body.Value
	return &topi.RequestBody{
		Description: b.Description,
		Required:    b.Required,
		Conetnt:     convertContent(b.Content),
	}
}

func convertContent(content openapi3.Content) []*topi.MediaTypeContent {
	ret := make([]*topi.MediaTypeContent, 0)
	for k, v := range content {
		c := &topi.MediaTypeContent{
			MediaType: k,
			Schema:    convertSchema(v.Schema),
		}
		ret = append(ret, c)
	}
	return ret
}

func convertResponses(responses openapi3.Responses) []*topi.Response {
	ret := make([]*topi.Response, 0)
	for status, response := range responses {
		r := &topi.Response{
			StatusCode:  status,
			Description: *response.Value.Description, // required
			Conetnt:     convertContent(response.Value.Content),
			Headers:     convertHeaders(response.Value.Headers),
		}
		ret = append(ret, r)
	}
	// sort to fix order because openapi3.Responses is map
	sort.Slice(ret, func(i, j int) bool { return ret[i].StatusCode < ret[j].StatusCode })
	return ret
}

func convertHeaders(headers openapi3.Headers) []*topi.Header {
	ret := make([]*topi.Header, 0)
	for k, v := range headers {
		p := &topi.Parameter{
			// name/in must not be specified
			Description: v.Value.Description,
			Required:    v.Value.Required,
			Deprecated:  v.Value.Deprecated,
			Schema:      convertSchema(v.Value.Schema),
		}
		h := &topi.Header{
			Name:      k,
			Parameter: p,
		}
		ret = append(ret, h)
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
	return ret
}

func convertSecurityRequirements(rs *openapi3.SecurityRequirements) []*topi.SecurityRequirement {
	// fixme: use instead if OpenAPI Object has Security Requirement Object
	if rs == nil {
		return nil
	}
	ret := make([]*topi.SecurityRequirement, 0)
	for _, r := range *rs {
		req := &topi.SecurityRequirement{
			Schemes: make([]*topi.SecurityRequirementScheme, 0),
		}
		for k, v := range r {
			s := &topi.SecurityRequirementScheme{
				Key:    k,
				Scopes: v,
			}
			req.Schemes = append(req.Schemes, s)
		}
		ret = append(ret, req)
	}
	return ret
}

func convertComponents(components *openapi3.Components) *topi.Components {
	return &topi.Components{
		SecuritySchemes: convertSecuritySchemes(components.SecuritySchemes),
	}
}

func convertSecuritySchemes(schemes openapi3.SecuritySchemes) []*topi.SecurityScheme {
	ret := make([]*topi.SecurityScheme, 0)
	for k, v := range schemes {
		s := &topi.SecurityScheme{
			Key:              k,
			Type:             v.Value.Type,
			Description:      v.Value.Description,
			Name:             v.Value.Name,
			In:               v.Value.In,
			Scheme:           v.Value.Scheme,
			BearerFormat:     v.Value.BearerFormat,
			OpenIdConnectUrl: v.Value.OpenIdConnectUrl,
			OAuthFlows:       convertOAuthFlows(v.Value.Flows),
		}
		ret = append(ret, s)
	}
	return ret
}

func convertOAuthFlows(flows *openapi3.OAuthFlows) *topi.OAtuhFlows {
	if flows == nil {
		return nil
	}
	return &topi.OAtuhFlows{
		ImplicitFlow:                         convertOAtuhFlow(flows.Implicit),
		ResourceOwnerPasswordCredentialsFlow: convertOAtuhFlow(flows.Password),
		ClientCredentialsFlow:                convertOAtuhFlow(flows.ClientCredentials),
		AuthorizatonCodeFlow:                 convertOAtuhFlow(flows.AuthorizationCode),
	}
}

func convertOAtuhFlow(flow *openapi3.OAuthFlow) *topi.OAuthFlow {
	if flow == nil {
		return nil
	}
	return &topi.OAuthFlow{
		AuthorizationURL: flow.AuthorizationURL,
		TokenURL:         flow.TokenURL,
		RefreshURL:       flow.RefreshURL,
		Scopes:           convertScopes(flow.Scopes),
	}
}

func convertScopes(scopes map[string]string) []*topi.Scope {
	ret := make([]*topi.Scope, 0)
	for k, v := range scopes {
		s := &topi.Scope{
			Name:   k,
			Detail: v,
		}
		ret = append(ret, s)
	}
	return ret
}
