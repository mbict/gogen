package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
	"github.com/mbict/gogen/extends/api"
	"github.com/mbict/gogen/extends/api/codegen"
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
	m.Service = s
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
		gogen.AsObject(m.Results.Type).Add(ut.Name(), ut.AttributeExpr)
	case gogen.DataType:
		gogen.AsObject(m.Results.Type).Add(resolveReturnName(v), gogen.NewAttribute(v))
	case func():
		dslengine.Execute(v, m.Results)
	default:
		dslengine.InvalidArgError("Expected a DataType or a name of a UserType", v)
	}
}

// resolveReturnName is a helper function to resolve if a return type name needs to plural or singular
// When a array data type is used we add (only once) an 's' as suffix and we try to find the type name
// `ArrayOf(User)` resolves to "Users" if the usertype its name is "User"
// `ArrayOf(Date)` resolves to "Dates" cause the user type name is "Date"
// `ArrayOf(ArrayOf(Int))` resolves to "Ints" nested arrays are ignored and will not add a extra 's'
// `User` resolves to "User"
// `Int` reolves to "Int"
func resolveReturnName(dt gogen.DataType) string {
	var f func(dt gogen.DataType) string
	suffix := ""

	f = func(dt gogen.DataType) string {
		if v, ok := dt.(*gogen.Array); ok {
			suffix = "s"
			return f(v.ElemType.Type)
		}
		return dt.Name()
	}

	return f(dt) + suffix
}

func Payload(dsl func()) {
	m, ok := methodDefinition()
	if !ok {
		return
	}

	dslengine.Execute(dsl, m.Payload)
}
