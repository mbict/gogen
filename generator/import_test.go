package generator

import (
	. "gopkg.in/check.v1"
)

type ImportSuite struct{}

var _ = Suite(&ImportSuite{})

func (s *ImportSuite) TestPackageName(c *C) {
	tests := []struct {
		path     string
		expected string
	}{
		{
			path:     "/a/b/test",
			expected: "test",
		}, {
			path:     "/a/b/test.blup",
			expected: "test.blup",
		}, {
			path:     "/a/b/go.test",
			expected: "test",
		}, {
			path:     "/a/b/go_test",
			expected: "test",
		}, {
			path:     "/a/b/go-test",
			expected: "test",
		}, {
			path:     "/a/b/go.go-test",
			expected: "go-test",
		}, {
			path:     "/a/b/go_go.test",
			expected: "go.test",
		}, {
			path:     "/a/b/go-go_test",
			expected: "go_test",
		},
	}

	for _, test := range tests {
		name := PackageName(test.path)

		c.Check(name, Equals, test.expected, Commentf("Failed output check path `%s` resolves to package name (expected) `%s` == `%s` (actual)", test.path, test.expected, name))
	}
}
