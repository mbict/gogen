package dsl

import (
	"github.com/mbict/gogen/dslengine"
	"github.com/mbict/gogen/extends/api"
)

func HTTP(dsl func()) {
	switch t := dslengine.Current().(type) {
	case *api.Service:
		dslengine.Execute(dsl, t.HttpEndpoint)
	case *api.Method:
		e := api.NewHTTPEndpoint()
		e.Method = t
		e.Parent = t.Service.HttpEndpoint
		dslengine.Execute(dsl, e)
		t.AddEndpoint(e)
	default:
		dslengine.IncompatibleDSL()
		return
	}
}

func Response(params ...interface{}) {

}

//Path sets the base path of a http endpoint
func Path(path string) {
	d, ok := httpDefinition()
	if !ok {
		return
	}
	d.SetPath(path)
}

func GET(path string) {
	d, ok := httpEndpointDefinition()
	if !ok {
		return
	}
	d.AddRoute("GET", path)
}

func POST(path string) {
	d, ok := httpEndpointDefinition()
	if !ok {
		return
	}
	d.AddRoute("POST", path)
}

func PUT(path string) {
	d, ok := httpEndpointDefinition()
	if !ok {
		return
	}
	d.AddRoute("PUT", path)
}

func DELETE(path string) {
	d, ok := httpEndpointDefinition()
	if !ok {
		return
	}
	d.AddRoute("DELETE", path)
}
