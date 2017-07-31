package dsl

import (
	"github.com/mbict/gogen/extends/api"
	"goa.design/goa.v2/eval"
	"github.com/mbict/gogen/extends/api/codegen"
	"github.com/mbict/gogen/dslengine"
)

// Service will describe a service endpoint inside the api
func Service(name string, dsl eval.DSLFunc) *api.Service {
	codegen.Register()

	if rootDefinition() == false {
		return nil
	}

	s := api.NewService(name)
	dslengine.Execute(dsl, s)
	api.Root.AddService(s)

	return s
}

// Method describes a method inside a service
func Method(name string, dsl eval.DSLFunc) *api.Method {
	s, ok := serviceDefinition()
	if !ok {
		return nil
	}

	m := api.NewMethod(name)
	s.AddMethod(m)

	return m
}

//Result describes the type that a service will return
//Provide a type a type or a dsl function with the attributes
func Result(typeDefintion interface{}) {

}

func Payload(dsl eval.DSLFunc) {

}
