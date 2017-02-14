package dslengine

// TopLevelDefinition represents the top-level file definitions, done
// with `var _ = `.  An instance of this object is returned by
// `CurrentDefinition()` when at the top-level.
type TopLevelDefinition struct{}

// Context tells the DSL engine which context we're in when showing
// errors.
func (t *TopLevelDefinition) Context() string {
	return "top-level"
}
