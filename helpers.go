package gogen

// Convenience methods

// AsObject returns the type underlying object if any, nil otherwise.
func AsObject(dt DataType) *Object {
	switch t := dt.(type) {
	case *UserTypeExpr:
		return AsObject(t.Type)
	case Composite:
		return AsObject(t.Attribute().Type)
	case *Object:
		return t
	default:
		return nil
	}
}

// AsArray returns the type underlying array if any, nil otherwise.
func AsArray(dt DataType) *Array {
	switch t := dt.(type) {
	case *UserTypeExpr:
		return AsArray(t.Type)
	case Composite:
		return AsArray(t.Attribute().Type)
	case *Array:
		return t
	default:
		return nil
	}
}

// AsMap returns the type underlying map if any, nil otherwise.
func AsMap(dt DataType) *Map {
	switch t := dt.(type) {
	case *UserTypeExpr:
		return AsMap(t.Type)
	case Composite:
		return AsMap(t.Attribute().Type)
	case *Map:
		return t
	default:
		return nil
	}
}

// IsObject returns true if the data type is an object.
func IsObject(dt DataType) bool {
	return AsObject(dt) != nil
}

// IsArray returns true if the data type is an array.
func IsArray(dt DataType) bool {
	return AsArray(dt) != nil
}

// IsMap returns true if the data type is a map.
func IsMap(dt DataType) bool {
	return AsMap(dt) != nil
}

// IsPrimitive returns true if the data type is a primitive type.
func IsPrimitive(dt DataType) bool {
	switch t := dt.(type) {
	case Primitive:
		return true
	case *UserTypeExpr:
		return IsPrimitive(t.Type)
	case Composite:
		return IsPrimitive(t.Attribute().Type)
	default:
		return false
	}
}

// sanitizeMapKeyName returns a trimmed name and a lowercase keyname
//func sanitzeMapKeyName(key string) (string, string) {
//	key = strings.TrimSpace(key)
//	return strings.ToLower(key), key
//}
