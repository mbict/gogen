package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

func ArrayOf(v interface{}, dsl ...func()) *gogen.Array {
	var t gogen.DataType
	var ok bool
	t, ok = v.(gogen.DataType)
	if !ok {
		if name, ok := v.(string); ok {
			t = gogen.Root.UserType(name)
		}
	}
	// never return nil to avoid panics, errors are reported after DSL execution
	res := &gogen.Array{ElemType: &gogen.AttributeExpr{Type: gogen.String}}
	if t == nil {
		dslengine.ReportError("invalid ArrayOf argument: not a type and not a known user type name")
		return res
	}
	if len(dsl) > 1 {
		dslengine.ReportError("ArrayOf: too many arguments")
		return res
	}
	at := gogen.AttributeExpr{Type: t}
	if len(dsl) == 1 {
		dslengine.Execute(dsl[0], &at)
	}
	return &gogen.Array{ElemType: &at}
}
