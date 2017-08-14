package api

import (
	"github.com/mbict/gogen/dslengine"
	"path"
	"regexp"
)

var WildcardRegex = regexp.MustCompile(`/(?::|\*)([a-zA-Z0-9_]+)`)

type PathResolver interface {
	SetPath(path string)
	Path() string
}

type Endpoint interface {
	dslengine.Definition
}

type HTTPRoute struct {
	Parent *HTTPEndpoint
	Method string
	path   string
}

func (r *HTTPRoute) Path() string {
	return path.Join(r.Parent.Path(), r.path)
}

func (r *HTTPRoute) SetPath(path string) {
	r.path = path
}

//Returns the names of the path params in this route
func (r *HTTPRoute) Params() []string {
	matches := WildcardRegex.FindAllStringSubmatch(path, -1)
	wcs := make([]string, len(matches))
	for i, m := range matches {
		wcs[i] = m[1]
	}
	return wcs
}

//HTTP is the base description for a http transport layer
type HTTP struct {
	Parent *HTTP
	path   string
}

func (*HTTP) Context() string {
	return "api.endpoint.http"
}

func (h *HTTP) SetPath(path string) {
	h.path = path
}

func (h *HTTP) Path() string {
	if h.Parent != nil {
		return path.Join(h.Parent.Path(), h.path)
	}
	return h.path
}

func NewHTTP() *HTTP {
	return &HTTP{}
}

//HttpEndpoint endpoint
type HTTPEndpoint struct {
	HTTP
	Routes []*HTTPRoute
	Method *Method
}

func (*HTTPEndpoint) GetParams() {

}

func (*HTTPEndpoint) GetQueryParams() {

}

func (*HTTPEndpoint) GetPathParams() {

}

func (*HTTPEndpoint) Context() string {
	return "api.endpoint.httpTransport"
}

func (h *HTTPEndpoint) AddRoute(method string, path string) {
	h.Routes = append(h.Routes, &HTTPRoute{
		Parent: h,
		Method: method,
		path:   path,
	})
}

func NewHTTPEndpoint() *HTTPEndpoint {
	return &HTTPEndpoint{}
}

//GRPC endpoint
type GRPC struct {
}

func (*GRPC) Context() string {
	return "api.endpoint.grpc"
}

func NewGRPCEndpoint() *GRPC {
	return &GRPC{}
}
