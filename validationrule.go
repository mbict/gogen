package gogen

type ValidationExpr struct {
	// Format represents a standardized format e.g. email,hostname, date-time etc
	//Format string
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
	// Required fields for object notations
	Required []string
}


func (v *ValidationExpr) AddRequired(required []string) {
	for _, r := range required {
		found := false
		for _, rr := range v.Required {
			if r == rr {
				found = true
				break
			}
		}
		if !found {
			v.Required = append(v.Required, r)
		}
	}
}