package dsl

import (
	"github.com/mbict/gogen/extends/api"
	"github.com/mbict/gogen/extends/api/codegen"
	"github.com/mbict/gogen/dslengine"
	"github.com/mbict/gogen"
	"fmt"
)

// Service will describe a service endpoint inside the api
func Service(name string, dsl func()) *api.Service {
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
func Method(name string, dsl func()) *api.Method {
	s, ok := serviceDefinition()
	if !ok {
		return nil
	}

	m := api.NewMethod(name)
	dslengine.Execute(dsl, m)
	s.AddMethod(m)

	return m
}

//Result describes the type that a service will return
//Provide a type a type or a dsl function with the attributes
func Result(t interface{}) {
	m, ok := methodDefinition()
	if !ok {
		return
	}

	switch v := t.(type) {
	case string:

		ut := gogen.Root.UserType(v)
		if ut == nil {
			dslengine.ReportError("unable to find result type for `%v`", t)
		}
		m.Results = append( m.Results, ut)
	case gogen.DataType:
		fmt.Println(v.Name(), v.Kind())
		m.Results = append( m.Results, v)
	case func():



	default:
		dslengine.InvalidArgError("", "")
		//dslengine.ReportError("unable to find result type for `%v`", t)
	}
}

func Payload(dsl func()) {

}
