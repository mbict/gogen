package gogen

type ValidationRule struct {
	// Pattern represents a regex pattern validation
	Pattern string
	// Minimum represents an minimum value validation
	Minimum *float64
	// Maximum represents a maximum value validation
	Maximum *float64
	// MinLength represents an minimum length validation
	MinLength *int
	// MaxLength represents an maximum length validation
	MaxLength *int
	// Required
	Required bool
}
