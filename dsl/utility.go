package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

// Package can be used on expressions who implement the dsl.Packager interface
// It is intended to describe the name of the package/namespace or directory where the generated code will be generated.
func Package(name string) {
	expr, ok := dslengine.Current().(gogen.Packager)
	if !ok {
		dslengine.IncompatibleDSL()
		return
	}
	expr.SetPackage(name)
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

func Error(name string, dataType gogen.DataType, description ...string) {

}

// Trait simply returns the dsl function for reuse
func Trait(dsl func()) func() {
	return dsl
}
