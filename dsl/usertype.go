package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

// Type defines a user type. A user type has a unique name and may be an alias
// to an existing type or may describe a completely new type using a list of
// attributes (object fields). Attribute types may themselves be user type.
// When a user type is defined as an alias to another type it may define
// additional validations - for example it a user type which is an alias of
// String may define a validation pattern that all instances of the type
// must match.
//
// Type is a top level definition.
//
// Type takes two or three arguments: the first argument is the name of the type.
// The name must be unique. The second argument is either another type or a
// function. If the second argument is a type then there may be a function passed
// as third argument.
//
// Example:
//
//     // simple alias
//     var MyString = Type("MyString", String)
//
//     // alias with description and additional validation
//     var Hostname = Type("Hostname", String, func() {
//         DescriptionExpr("A host name")
//         Format(FormatHostname)
//     })
//
//     // new type
//     var SumPayload = Type("SumPayload", func() {
//         DescriptionExpr("Type sent to add endpoint")
//
//         Attribute("a", String)                 // string attribute "a"
//         Attribute("b", Int32, "operand")       // attribute with description
//         Attribute("operands", ArrayOf(Int32))  // array attribute
//         Attribute("ops", MapOf(String, Int32)) // map attribute
//         Attribute("c", SumMod)                 // attribute using user type
//         Attribute("len", Int64, func() {       // attribute with validation
//             Minimum(1)
//         })
//
//         Required("a")                          // Required attributes
//         Required("b", "c")
//     })
//
func Type(name string, args ...interface{}) *gogen.UserTypeExpr {
	if len(args) > 2 {
		dslengine.ReportError("too many arguments")
		return nil
	}
	if t := gogen.Root.UserType(name); t != nil {
		dslengine.ReportError("type %#v defined twice", name)
		return nil
	}

	if _, ok := dslengine.Current().(*dslengine.TopLevelDefinition); !ok {
		dslengine.IncompatibleDSL()
		return nil
	}

	var (
		base gogen.DataType
		dsl  func()
	)
	if len(args) == 0 {
		// Make Type behave like Attribute
		args = []interface{}{gogen.String}
	}
	switch a := args[0].(type) {
	case gogen.DataType:
		base = a
		if len(args) == 2 {
			d, ok := args[1].(func())
			if !ok {
				dslengine.ReportError("third argument must be a function")
				return nil
			}
			dsl = d
		}
	case func():
		base = &gogen.Object{}
		dsl = a
		if len(args) == 2 {
			dslengine.ReportError("only one argument allowed when it is a function")
			return nil
		}
	default:
		dslengine.InvalidArgError("type or function", args[0])
		return nil
	}

	ut := &gogen.UserTypeExpr{
		TypeName:      name,
		AttributeExpr: gogen.NewAttribute(base),
	}
	dslengine.Execute(dsl, ut)

	gogen.Root.AddUserType(ut)
	return ut
}
