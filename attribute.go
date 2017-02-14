package gogen

// Attribute defines the base of attribute
type AttributeExpr struct {
	// Type is attribute data type e.g. Primitive, array, object
	Type DataType

	// Vaidation describes the rules for validation of a attribute
	Validation *ValidationRule

	// Metadata is a list of Key/Values pair
	Metadata MetadataList
}

func (a *AttributeExpr) Context() string {
	return "attribute"
}
