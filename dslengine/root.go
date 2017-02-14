package dslengine

// Root is the interface implemented by the DSL root objects.
// These objects contains all the definition sets created by the DSL and can
// be passed to the dsl engine for execution.
type Root interface {

	// DSLName is displayed by the runner upon executing the DSL.
	// Registered DSL roots must have unique names.
	DSLName() string

	// IterateSets implements the visitor pattern: is is called by the engine so the
	// DSL can control the order of execution. IterateSets calls back the engine via
	// the given iterator as many times as needed providing the DSL definitions that
	// must be run for each callback.
	IterateSets(SetIterator)

	// Reset restores the root to pre DSL execution state.
	// This is mainly used by tests.
	Reset()
}
