package dslengine

// DSL evaluation contexts stack
type context []Definition

// Current evaluation context, i.e. object being currently built by DSL
func (s context) Current() Definition {
	if len(s) == 0 {
		return nil
	}
	return s[len(s)-1]
}

// Reset context stack to root
func (s context) Reset() {
	s = s[:0]
}
