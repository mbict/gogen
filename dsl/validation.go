package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
	"reflect"
	"regexp"
	"strconv"
)

func Required(names ...string) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil && a.Type.Kind() != gogen.ObjectKind {
		incompatibleAttributeType("required", a.Type.Name(), "an object")
	} else {
		if a.Validation == nil {
			a.Validation = &gogen.ValidationExpr{}
		}
		a.Validation.AddRequired(names)
	}
}

//
//func IsRequired() {
//	a := attributeFromContext()
//	if a == nil {
//		dslengine.IncompatibleDSL()
//		return
//	}
//
//	a.Validation.Required = true
//}

func Pattern(p string) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil && a.Type.Kind() != gogen.StringKind {
		incompatibleAttributeType("pattern", a.Type.Name(), "a string")
	} else {
		_, err := regexp.Compile(p)
		if err != nil {
			dslengine.ReportError("invalid pattern %#v, %s", p, err)
		} else {
			if a.Validation == nil {
				a.Validation = &gogen.ValidationExpr{}
			}
			a.Validation.Pattern = p
		}
	}
}

func Minimum(val interface{}) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil &&
		a.Type.Kind() != gogen.Int32Kind && a.Type.Kind() != gogen.Int64Kind &&
		a.Type.Kind() != gogen.Float32Kind && a.Type.Kind() != gogen.Float64Kind {

		incompatibleAttributeType("minimum", a.Type.Name(), "an integer or a number")
	} else {
		var f float64
		switch v := val.(type) {
		case float32, float64, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			f = reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0.0))).Float()
		case string:
			var err error
			f, err = strconv.ParseFloat(v, 64)
			if err != nil {
				dslengine.ReportError("invalid number value %#v", v)
				return
			}
		default:
			dslengine.ReportError("invalid number value %#v", v)
			return
		}

		if a.Validation == nil {
			a.Validation = &gogen.ValidationExpr{}
		}
		a.Validation.Minimum = &f
	}
}

func Maximum(val interface{}) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil &&
		a.Type.Kind() != gogen.Int32Kind && a.Type.Kind() != gogen.Int64Kind &&
		a.Type.Kind() != gogen.Float32Kind && a.Type.Kind() != gogen.Float64Kind {

		incompatibleAttributeType("maximum", a.Type.Name(), "an integer or a number")
	} else {
		var f float64
		switch v := val.(type) {
		case float32, float64, int, int8, int16, int32, int64, uint8, uint16, uint32, uint64:
			f = reflect.ValueOf(v).Convert(reflect.TypeOf(float64(0.0))).Float()
		case string:
			var err error
			f, err = strconv.ParseFloat(v, 64)
			if err != nil {
				dslengine.ReportError("invalid number value %#v", v)
				return
			}
		default:
			dslengine.ReportError("invalid number value %#v", v)
			return
		}

		if a.Validation == nil {
			a.Validation = &gogen.ValidationExpr{}
		}
		a.Validation.Maximum = &f
	}
}

func MinLength(val int) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil {
		kind := a.Type.Kind()
		if kind != gogen.BytesKind &&
			kind != gogen.StringKind &&
			kind != gogen.ArrayKind &&
			kind != gogen.MapKind {

			incompatibleAttributeType("minimum length", a.Type.Name(), "a string or an array")
			return
		}
	}

	if a.Validation == nil {
		a.Validation = &gogen.ValidationExpr{}
	}
	a.Validation.MinLength = &val
}

func MaxLength(val int) {
	a := attributeFromContext()
	if a == nil {
		dslengine.IncompatibleDSL()
		return
	}

	if a.Type != nil {
		kind := a.Type.Kind()
		if kind != gogen.BytesKind &&
			kind != gogen.StringKind &&
			kind != gogen.ArrayKind &&
			kind != gogen.MapKind {

			incompatibleAttributeType("maximum length", a.Type.Name(), "a string or an array")
			return
		}
	}

	if a.Validation == nil {
		a.Validation = &gogen.ValidationExpr{}
	}
	a.Validation.MaxLength = &val
}

// incompatibleAttributeType reports an error for validations defined on
// incompatible attributes (e.g. max value on string).
func incompatibleAttributeType(validation, actual, expected string) {
	dslengine.ReportError("invalid %s validation definition: attribute must be %s (but type is %s)", validation, expected, actual)
}

func attributeFromContext() *gogen.AttributeExpr {
	switch def := dslengine.Current().(type) {
	case *gogen.AttributeExpr:
		return def
	case gogen.Composite:
		return def.Attribute()
	}
	return nil
}
