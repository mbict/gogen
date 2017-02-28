package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

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

	t := &gogen.UserTypeExpr{
		TypeName:      name,
		AttributeExpr: &gogen.AttributeExpr{Type: base},
	}
	dslengine.Execute(dsl, t.AttributeExpr)

	//gogen.Root.Types = append(gogen.Root.Types, t)
	return t
}
