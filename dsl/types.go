package dsl

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/dslengine"
)

//alias for types
const (
	Boolean = gogen.Boolean
	Int32   = gogen.Int32
	Int64   = gogen.Int64
	UInt32  = gogen.UInt32
	UInt64  = gogen.UInt64
	Float32 = gogen.Float32
	Float64 = gogen.Float64
	String  = gogen.String
	Bytes   = gogen.Bytes
	Any     = gogen.Any
)

// Usertype of predefined UUID type
var UUID = gogen.NewUserType("UUID", String, "/github.com/satori/go.uuid")

// Usertype of the date / time format type
var DateTime = gogen.NewUserType("Time", String, "/time")

// Usertype of the go builtin context
var Context = gogen.NewUserType("Context", Any, "/context")

// ArrayOf creates an array type from its element type.
//
// ArrayOf may be used wherever types can.
// The first argument of ArrayOf is the type of the array elements specified by
// name or by reference.
// The second argument of ArrayOf is an optional function that defines
// validations for the array elements.
//
// Examples:
//
//    var Names = ArrayOf(String, func() {
//        Pattern("[a-zA-Z]+") // Validates elements of the array
//    })
//
//    var Account = Type("Account", func() {
//        Attribute("bottles", ArrayOf(Bottle), "Account bottles", func() {
//            MinLength(1) // Validates array as a whole
//        })
//    })
//
// Note: CollectionOf and ArrayOf both return array types. CollectionOf returns
// a media type where ArrayOf returns a user type. In general you want to use
// CollectionOf if the argument is a media type and ArrayOf if it is a user
// type.
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
	res := &gogen.Array{ElemType: gogen.NewAttribute(gogen.String)}
	if t == nil {
		dslengine.ReportError("invalid ArrayOf argument: not a type and not a known user type name")
		return res
	}
	if len(dsl) > 1 {
		dslengine.ReportError("ArrayOf: too many arguments")
		return res
	}
	at := gogen.NewAttribute(t)
	if len(dsl) == 1 {
		dslengine.Execute(dsl[0], at)
	}
	return &gogen.Array{ElemType: at}
}

// MapOf creates a map from its key and element types.
//
// MapOf may be used wherever types can.
// MapOf takes two arguments: the key and value types either by name of by reference.
//
// Example:
//
//    var ReviewByID = MapOf(Int64, String, func() {
//        Key(func() {
//            Minimum(1)           // Validates keys of the map
//        })
//        Value(func() {
//            Pattern("[a-zA-Z]+") // Validates values of the map
//        })
//    })
//
//    var Review = Type("Review", func() {
//        Attribute("ratings", MapOf(Bottle, Int32), "Bottle ratings")
//    })
//
//func MapOf(k, v interface{}, dsl ...func()) *gogen.Map {
//	var tk, tv gogen.DataType
//	var ok bool
//	tk, ok = k.(gogen.DataType)
//	if !ok {
//		if name, ok := k.(string); ok {
//			tk = gogen.Root.UserType(name)
//		}
//	}
//	tv, ok = v.(gogen.DataType)
//	if !ok {
//		if name, ok := v.(string); ok {
//			tv = gogen.Root.UserType(name)
//		}
//	}
//	// never return nil to avoid panics, errors are reported after DSL execution
//	res := &gogen.Map{KeyType: &gogen.AttributeExpr{Type: gogen.String}, ElemType: &gogen.AttributeExpr{Type: gogen.String}}
//	if tk == nil {
//		dslengine.ReportError("invalid MapOf key argument: not a type and not a known user type name")
//		return res
//	}
//	if tv == nil {
//		dslengine.ReportError("invalid MapOf value argument: not a type and not a known user type name")
//		return res
//	}
//	if len(dsl) > 1 {
//		dslengine.ReportError("MapOf: too many arguments")
//		return res
//	}
//	kat := gogen.AttributeExpr{Type: tk}
//	vat := gogen.AttributeExpr{Type: tv}
//	m := &gogen.Map{KeyType: &kat, ElemType: &vat}
//	if len(dsl) == 1 {
//		mat := gogen.AttributeExpr{Type: m}
//		dslengine.Execute(dsl[0], &mat)
//	}
//	return m
//}
//
//// Key makes it possible to specify validations for map keys.
//func Key(dsl func()) {
//	at, ok := dslengine.Current().(*gogen.AttributeExpr)
//	if !ok {
//		dslengine.IncompatibleDSL()
//		return
//	}
//	m := at.Type.(*gogen.Map)
//	dslengine.Execute(dsl, m.KeyType)
//}
//
//// Value makes it possible to specify validations for map values.
//func Value(dsl func()) {
//	at, ok := dslengine.Current().(*gogen.AttributeExpr)
//	if !ok {
//		dslengine.IncompatibleDSL()
//		return
//	}
//	m := at.Type.(*gogen.Map)
//	dslengine.Execute(dsl, m.ElemType)
//}
