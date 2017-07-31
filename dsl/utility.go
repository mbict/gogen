package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

// Namespace can be used on expressions who implement the dsl.Namespace interface
// It is intended to describe the name of the package/namespace or directory where the generated code will be generated.
func Namespace(namespace string) {
	expr, ok := dslengine.Current().(gogen.Namespace)
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}
	expr.SetNamespace(namespace)
}

// DescriptionExpr will set the docuementation/description about a particular expression.
// This dsl can be used on expressions who implement the dsl.Describer interface.
func Description(description string) {
	expr, ok := dslengine.Current().(gogen.Describer)
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}
	expr.SetDescription(description)
}

func Error( name string, dataType gogen.DataType, description ...string) {

}
