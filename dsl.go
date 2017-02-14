package gogen

import (
	"github.com/dimfeld/httppath"
	"goa.design/goa.v2/goa/eval"
	"path"
	"strings"
)

// Endpoint defines a single endpoint / handler
type Endpoint struct {
	// Name of the endpoint
	Name string

	// Description about the endpoint
	Description string

	// Request defines the input params
	Request *UserTypeExpr

	// Respons defines the response the endpoint returns
	Response *UserTypeExpr
}

type Root struct {
	// Path is the base path prefix to all the endpoints.
	Path string

	// Resources contains the resources
	Resources []*Resource
}

type Resource struct {
	// Name of the resource
	Name string

	// Path is the URL path e.g. "/tasks/{id}"
	Path string

	// ParentName is the name of the parent resource
	ParentName string

	// Action with canonical resource path
	CanonicalActionName string

	// Actions contains the actions for this route
	Actions []*Action

	//The base root
	Root *Root
}

type Action struct {
	*Endpoint

	// Routes
	Routes []*Route

	// Resource is the parent resource this action belongs to
	Resource *Resource

	// Path and query string parameters
	Params *AttributeExpr

	// Payload is the payload
	Payload *AttributeExpr

	// Responses
	Responses []*HttpResponse
}

type HttpResponse struct {
	// Response name
	Name string
	// HTTP status
	Status int
	// Response description
	Description string
	// Response body type if any
	Type DataType
	// Response body media type if any
	MediaType string
	// Response view name if MediaType is the id of a MediaTypeExpr
	ViewName string
	// Response header expressions
	Headers *AttributeExpr
	// Parent action or resource
	Parent eval.Expression
	// Standard is true if the response is one of the default responses.
	Standard bool
}

type Route struct {
	// Method is the HTTP method e.g. GET, POST, PUT, PATCH, DELETE
	Method string

	// Path is the URL path e.g. "/tasks/{id}"
	Path string

	// Action is the parent action this route belongs to
	Action *Action
}

//------------------
// Root
//------------------

// Resource returns the resource with the given name if any.
func (r *Root) Resource(name string) *Resource {
	for _, res := range r.Resources {
		if res.Name == name {
			return res
		}
	}
	return nil
}

//------------------
// Resource
//------------------
// FullPath computes the base path to the resource actions concatenating the API and parent resource
// base paths as needed.
func (r *Resource) FullPath() string {
	if strings.HasPrefix(r.Path, "//") {
		return httppath.Clean(r.Path)
	}
	var basePath string
	if p := r.Parent(); p != nil {
		if ca := p.CanonicalAction(); ca != nil {
			if routes := ca.Routes; len(routes) > 0 {
				basePath = path.Join(routes[0].FullPath())
			}
		}
	} else {
		basePath = r.Root.Path
	}
	return httppath.Clean(path.Join(basePath, r.Path))
}

// CanonicalAction returns the canonical action of the resource if any.
// The canonical action is used to compute the default url to resources.
func (r *Resource) CanonicalAction() *Action {
	name := r.CanonicalActionName
	if name == "" {
		name = "show"
	}
	return r.Action(name)
}

// Action returns the resource action with the given name or nil if there isn't one.
func (r *Resource) Action(name string) *Action {
	for _, a := range r.Actions {
		if a.Name == name {
			return a
		}
	}
	return nil
}

// Parent returns the parent resource, nil otherwise.
func (r *Resource) Parent() *Resource {
	if r.ParentName != "" {
		if parent := r.Root.Resource(r.ParentName); parent != nil {
			return parent
		}
	}
	return nil
}

//------------------
// Action
//------------------

//------------------
// Route
//------------------

// Params returns the route parameters.
// For example for the route "GET /foo/:fooID" Params returns []string{"fooID"}.
//func (r *Route) Params() []string {
//	return ExtractRouteWildcards(r.FullPath())
//}

// FullPath returns the action full path computed by concatenating the API and resource base paths
// with the action specific path.
func (r *Route) FullPath() string {
	if r.IsAbsolute() {
		return httppath.Clean(r.Path[1:])
	}
	var base string
	if r.Action != nil && r.Action.Resource != nil {
		base = r.Action.Resource.FullPath()
	}
	return httppath.Clean(path.Join(base, r.Path))
}

// IsAbsolute returns true if the action path is a absolute path starting with a "//"
func (r *Route) IsAbsolute() bool {
	return strings.HasPrefix(r.Path, "//")
}
