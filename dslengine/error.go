package dslengine

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// Error represents an error that occurred while running the API DSL.
// It contains the name of the file and line number of where the error
// occurred as well as the original Go error.
type Error struct {
	GoError error
	File    string
	Line    int
}

// MultiError collects all DSL errors. It implements error.
type MultiError []*Error

// Error returns the error message.
func (m MultiError) Error() string {
	msgs := make([]string, len(m))
	for i, de := range m {
		msgs[i] = de.Error()
	}
	return strings.Join(msgs, "\n")
}

// Error returns the underlying error message.
func (de *Error) Error() string {
	if err := de.GoError; err != nil {
		if de.File == "" {
			return err.Error()
		}
		return fmt.Sprintf("[%s:%d] %s", de.File, de.Line, err.Error())
	}
	return ""
}

// computeErrorLocation implements a heuristic to find the location in the user
// code where the error occurred. It walks back the callstack until the file
// doesn't match "/goa/design/*.go" or one of the DSL package paths.
// When successful it returns the file name and line number, empty string and
// 0 otherwise.
func computeErrorLocation() (file string, line int) {
	skipFunc := func(file string) bool {
		if strings.HasSuffix(file, "_test.go") {
			// Be nice with tests
			return false
		}
		file = filepath.ToSlash(file)
		for pkg := range dslPackages {
			if strings.Contains(file, pkg) {
				return true
			}
		}
		return false
	}
	depth := 2
	_, file, line, _ = runtime.Caller(depth)
	for skipFunc(file) {
		depth++
		_, file, line, _ = runtime.Caller(depth)
	}
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	wd, err = filepath.Abs(wd)
	if err != nil {
		return
	}
	f, err := filepath.Rel(wd, file)
	if err != nil {
		return
	}
	file = f
	return
}

// caller returns the name of calling function.
func caller() string {
	pc, file, _, ok := runtime.Caller(2)
	if ok && filepath.Base(file) == "current.go" {
		pc, _, _, ok = runtime.Caller(3)
	}
	if !ok {
		return "<unknown>"
	}

	return runtime.FuncForPC(pc).Name()
}
