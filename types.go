package gogen

import "errors"

type Kind int

type DataType interface {
	// Kind of data type, one of the Kind enum.
	Kind() Kind
	// Name returns the type name.
	Name() string
}

type Primitive Kind

//type Object map[string]*Attribute

type Field struct {
	Name      string
	Attribute *AttributeExpr
}

// Object is used to describe composite data structures.
type Object []*Field

// Array is the type used to describe field arrays or repeated fields.
type Array struct {
	ElemType *AttributeExpr
}

// Map is the type used to describe maps of fields.
type Map struct {
	KeyType  *AttributeExpr
	ElemType *AttributeExpr
}

// Composite interface embeds the attribute definition
type Composite interface {
	// Attribute returns the composite expression embedded attribute.
	Attribute() *AttributeExpr
}

// UserTypeExpr defines a user defined type
type UserTypeExpr struct {
	*AttributeExpr

	// TypeName defines the name of the new user type
	TypeName string

	// Package this user type belongs to
	Package string
}

const (
	// BooleanKind represents a boolean.
	BooleanKind Kind = iota + 1
	// Int32Kind represents a signed 32-bit integer.
	Int32Kind
	// Int64Kind represents a signed 64-bit integer.
	Int64Kind
	// UInt32Kind represents an unsigned 32-bit integer.
	UInt32Kind
	// UInt64Kind represents an unsigned 64-bit integer.
	UInt64Kind
	// Float32Kind represents a 32-bit floating number.
	Float32Kind
	// Float64Kind represents a 64-bit floating number.
	Float64Kind
	// StringKind represents a JSON string.
	StringKind
	// BytesKind represent a series of bytes (binary data).
	BytesKind
	// ArrayKind represents a JSON array.
	ArrayKind
	// ObjectKind represents a JSON object.
	ObjectKind
	// MapKind represents a JSON object where the keys are not known in advance.
	MapKind
	// UserTypeKind represents a user type.
	UserTypeKind
	// MediaTypeKind represents a media type.
	MediaTypeKind
	// AnyKind represents a unknown type.
	AnyKind
)

const (
	// Boolean is the type for a JSON boolean.
	Boolean = Primitive(BooleanKind)

	// Int32 is the type for a signed 32-bit integer.
	Int32 = Primitive(Int32Kind)

	// Int64 is the type for a signed 64-bit integer.
	Int64 = Primitive(Int64Kind)

	// UInt32 is the type for an unsigned 32-bit integer.
	UInt32 = Primitive(UInt32Kind)

	// UInt64 is the type for an unsigned 64-bit integer.
	UInt64 = Primitive(UInt64Kind)

	// Float32 is the type for a 32-bit floating number.
	Float32 = Primitive(Float32Kind)

	// Float64 is the type for a 64-bit floating number.
	Float64 = Primitive(Float64Kind)

	// String is the type for a JSON string.
	String = Primitive(StringKind)

	// Bytes is the type for binary data.
	Bytes = Primitive(BytesKind)

	// Any is the type for an arbitrary JSON value (interface{} in Go).
	Any = Primitive(AnyKind)
)

// Empty represents empty values.
var Empty = &UserTypeExpr{
	TypeName: "Empty",
	AttributeExpr: &AttributeExpr{
		Type: &Object{},
	},
}

// Kind implements the DataKind interface
func (p Primitive) Kind() Kind {
	return Kind(p)
}

// Name returns the type name appropriate for logging.
func (p Primitive) Name() string {
	switch p {
	case Boolean:
		return "boolean"
	case Int32:
		return "int32"
	case Int64:
		return "int64"
	case UInt32:
		return "uint32"
	case UInt64:
		return "uint64"
	case Float32:
		return "float32"
	case Float64:
		return "float64"
	case String:
		return "string"
	case Bytes:
		return "[]byte"
	case Any:
		return "any"
	default:
		panic("unknown primitive type")
	}
}

// Kind implements DataKind interface
func (a *Array) Kind() Kind {
	return ArrayKind
}

// Name returns the type name.
func (a *Array) Name() string {
	return "array"
}

// Kind implements DataKind.
func (m *Map) Kind() Kind {
	return MapKind
}

// Name returns the type name.
func (m *Map) Name() string {
	return "hash"
}

// Kind implements DataKind.
func (o *Object) Kind() Kind {
	return ObjectKind
}

// Name returns the type name.
func (o *Object) Name() string {
	return "object"
}

func (o *Object) Add(name string, attr *AttributeExpr) error {
	if o.Field(name) != nil {
		return errors.New("duplicate field")
	}

	*o = append(*o, &Field{
		Name:      name,
		Attribute: attr,
	})
	return nil
}

func (o *Object) Field(name string) *AttributeExpr {
	for _, f := range *o {
		if f.Name == name {
			return f.Attribute
		}
	}
	return nil
}

// Kind implements DataKind.
func (ut *UserTypeExpr) Kind() Kind {
	return UserTypeKind
}

// Name returns the type name.
func (ut *UserTypeExpr) Name() string {
	return ut.TypeName
}

// Usertype of predefined UUID type
var UUID = &UserTypeExpr{
	AttributeExpr: &AttributeExpr{
		Type: String,
	},
	TypeName: "UUID",
	Package:  "github.com/satori/go.uuid",
}

// Usertype of the date / time format type
var Date = &UserTypeExpr{
	AttributeExpr: &AttributeExpr{
		Type: String,
	},
	TypeName: "Time",
	Package:  "time",
}
