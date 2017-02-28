package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
	"strings"
)

// Attribute describes a field of an object.
//
// An attribute has a name, a type and optionally a default value, an example
// value and validation rules.
//
// The type of an attribute can be one of:
//
// * The primitive types Boolean, Float32, Float64, Int, Int32, Int64, UInt,
//   UInt32, UInt64, String or Bytes.
//
// * A user type defined via the Type function.
//
// * An array defined using the ArrayOf function.
//
// * An map defined using the MapOf function.
//
// * An object defined inline using Attribute to define the type fields
//   recursively.
//
// * The special type Any to indicate that the attribute may take any of the
//   types listed above.
//
// Attribute may appear in MediaType, Type, Attribute or Attributes.
//
// Attribute accepts one to four arguments, the valid usages of the function
// are:
//
//    Attribute(name)       // Attribute of type String with no description, no
//                          // validation, default or example value
//
//    Attribute(name, dsl)  // Attribute of type object with inline field
//                          // definitions, description, validations, default
//                          // and/or example value
//
//    Attribute(name, type) // Attribute with no description, no validation,
//                          // no default or example value
//
//    Attribute(name, type, dsl) // Attribute with description, validations,
//                               // default and/or example value
//
//    Attribute(name, type, description)      // Attribute with no validation,
//                                            // default or example value
//
//    Attribute(name, type, description, dsl) // Attribute with description,
//                                            // validations, default and/or
//                                            // example value
//
// Where name is a string indicating the name of the attribute, type specifies
// the attribute type (see above for the possible values), description a string
// providing a human description of the attribute and dsl the defining DSL if
// any.
//
// When defining the type inline using Attribute recursively the function takes
// the second form (name and DSL defining the type). The description can be
// provided using the Description function in this case.
//
// Examples:
//
//    Attribute("name")
//
//    Attribute("driver", Person)         // Use type defined with Type function
//
//    Attribute("driver", "Person")       // May also use the type name
//
//    Attribute("name", String, func() {
//        Pattern("^foo")                 // Adds a validation rule
//    })
//
//    Attribute("driver", Person, func() {
//        Required("name")                // Add required field to list of
//    })                                  // fields already required in Person
//
//    Attribute("name", String, func() {
//        Default("bob")                  // Sets a default value
//    })
//
//    Attribute("name", String, "name of driver") // Sets a description
//
//    Attribute("age", Int32, "description", func() {
//        Minimum(2)                       // Sets both a description and
//                                         // validations
//    })
//
// The definition below defines an attribute inline. The resulting type
// is an object with three attributes "name", "age" and "child". The "child"
// attribute is itself defined inline and has one child attribute "name".
//
//    Attribute("driver", func() {           // Define type inline
//        Description("Composite attribute") // Set description
//
//        Attribute("name", String)          // Child attribute
//        Attribute("age", Int32, func() {   // Another child attribute
//            Description("Age of driver")
//            Default(42)
//            Minimum(2)
//        })
//        Attribute("child", func() {        // Defines a child attribute
//            Attribute("name", String)      // Grand-child attribute
//            Required("name")
//        })
//
//        Required("name", "age")            // List required attributes
//    })
//
func Attribute(name string, args ...interface{}) {
	name = strings.TrimSpace(name)
	var parent *gogen.AttributeExpr

	switch def := dslengine.Current().(type) {
	case *gogen.AttributeExpr:
		parent = def
	case gogen.Composite:
		parent = def.Attribute()
	default:
		dslengine.IncompatibleDSL()
		return
	}

	if parent != nil {
		if parent.Type == nil {
			parent.Type = &gogen.Object{}
		}
		if _, ok := parent.Type.(*gogen.Object); !ok {
			dslengine.ReportError("can't define child attributes on attribute of type %s", parent.Type.Name())
			return
		}

		var baseAttr *gogen.AttributeExpr
		//if parent.Reference != nil {
		//	if att, ok := gogen.AsObject(parent.Reference)[name]; ok {
		//		baseAttr = gogen.DupAtt(att)
		//	}
		//}

		dataType, description, dsl := parseAttributeArgs(baseAttr, args...)
		//if baseAttr != nil {
		//if description != "" {
		//	baseAttr.Description = description
		//}
		//	if dataType != nil {
		//		baseAttr.Type = dataType
		//	}
		//} else {
		baseAttr = &gogen.AttributeExpr{
			Type:        dataType,
			Description: description,
		}
		//}
		//baseAttr.Reference = parent.Reference
		if dsl != nil {
			dslengine.Execute(dsl, baseAttr)
		}

		//if no attribute type was provided we default back to string
		if baseAttr.Type == nil {
			baseAttr.Type = gogen.String
		}

		gogen.AsObject(parent.Type).Add(name, baseAttr)
	}
}

func parseAttributeArgs(baseAttr *gogen.AttributeExpr, args ...interface{}) (gogen.DataType, string, func()) {
	var (
		dataType    gogen.DataType
		description string
		dsl         func()
		ok          bool
	)

	parseDataType := func(expected string, index int) {
		if name, ok2 := args[index].(string); ok2 {
			// Lookup type by name
			if dataType = gogen.Root.UserType(name); dataType == nil {
				dslengine.InvalidArgError(expected, args[index])
			}
			return
		}
		if dataType, ok = args[index].(gogen.DataType); !ok {
			dslengine.InvalidArgError(expected, args[index])
		}
	}
	parseDescription := func(expected string, index int) {
		if description, ok = args[index].(string); !ok {
			dslengine.InvalidArgError(expected, args[index])
		}
	}
	parseDSL := func(index int, success, failure func()) {
		if dsl, ok = args[index].(func()); ok {
			success()
			return
		}
		failure()
	}

	success := func() {}

	switch len(args) {
	case 0:
		if baseAttr != nil {
			dataType = baseAttr.Type
		} else {
			dataType = gogen.String
		}
	case 1:
		success = func() {
			if baseAttr != nil {
				dataType = baseAttr.Type
			}
		}
		parseDSL(0, success, func() { parseDataType("type, type name or func()", 0) })
	case 2:
		parseDataType("type or type name", 0)
		parseDSL(1, success, func() { parseDescription("string or func()", 1) })
	case 3:
		parseDataType("type or type name", 0)
		parseDescription("string", 1)
		parseDSL(2, success, func() { dslengine.InvalidArgError("func()", args[2]) })
	default:
		dslengine.ReportError("too many arguments in call to Attribute")
	}

	return dataType, description, dsl
}
