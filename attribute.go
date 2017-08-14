package gogen

// Attribute defines the base of attribute
type AttributeExpr struct {
	// DescriptionExpr is a description about the attribute
	description string

	// Type is attribute data type e.g. Primitive, array, object
	Type DataType

	// Vaidation describes the rules for validation of a attribute
	Validation *ValidationExpr

	// Metadata is a list of Key/Values pair
	Metadata Metadata
}

func (a *AttributeExpr) Description() string {
	return a.description
}

func (a *AttributeExpr) SetDescription(description string) {
	a.description = description
}

func (a *AttributeExpr) Context() string {
	return "attribute"
}

// AllRequired returns the list of all required fields from the underlying
// object. This method recurses if the type is itself an attribute (i.e. a
// UserType, this happens with the Reference DSL for example).
func (a *AttributeExpr) AllRequired() (required []string) {
	if a == nil || a.Validation == nil {
		return
	}
	required = a.Validation.Required
	if u, ok := a.Type.(*UserTypeExpr); ok {
		required = append(required, u.Attribute().AllRequired()...)
	}
	return
}

// IsRequired returns true if the given string matches the name of a required
// attribute, false otherwise. This method only applies to attributes of type
// Object.
func (a *AttributeExpr) IsRequired(attributeName string) bool {
	for _, name := range a.AllRequired() {
		if name == attributeName {
			return true
		}
	}
	return false
}

func NewAttribute(dataType DataType) *AttributeExpr {
	return &AttributeExpr{
		Type:       dataType,
		Metadata:   make(Metadata, 0),
		Validation: &ValidationExpr{},
	}
}
