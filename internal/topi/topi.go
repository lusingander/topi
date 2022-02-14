package topi

type Document struct {
	TagPathMap map[string][]*Path
	Tags       []*Tag
}

type Path struct {
	UriPath     string
	Method      string
	OperationId string
	Summary     string
	Deprecated  bool
}

type Tag struct {
	Name        string
	Description string
}
