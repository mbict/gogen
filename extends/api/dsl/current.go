package dsl

import (
	"github.com/mbict/gogen/dslengine"
	"github.com/mbict/gogen/extends/api"
)

func rootDefinition() bool {
	_, ok := dslengine.Current().(*dslengine.TopLevelDefinition)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return ok
}

func serviceDefinition() (*api.Service, bool) {
	d, ok := dslengine.Current().(*api.Service)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}

func methodDefinition() (*api.Method, bool) {
	d, ok := dslengine.Current().(*api.Method)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}

func httpEndpointDefinition() (*api.HTTPEndpoint, bool) {
	d, ok := dslengine.Current().(*api.HTTPEndpoint)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}

func httpDefinition() (*api.HTTP, bool) {
	d, ok := dslengine.Current().(*api.HTTP)
	if !ok {
		dslengine.IncompatibleDSL()
	}
	return d, ok
}
