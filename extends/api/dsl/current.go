package dsl

import (
	"github.com/mbict/gogen/extends/api"
	"github.com/mbict/gogen/dslengine"
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
