package codegen

import (
	"bytes"
	"github.com/mbict/gogen"
	. "gopkg.in/check.v1"
	"io"
	"os"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type GeneratorSuite struct {
	generator *Codegen

	root *gogen.Gogen
}

var _ = Suite(&GeneratorSuite{})

func (s *GeneratorSuite) SetUpTest(c *C) {
	s.generator = NewCodeGenerator("testing/base")
	s.root = &gogen.Gogen{}
}

func (s *GeneratorSuite) TestGenerateType(c *C) {

	tests := []struct {
		name     string
		userType *gogen.UserTypeExpr
		expected string
	}{}

	for _, test := range tests {
		sec, err := s.generator.GenerateType(test.userType)

		c.Check(err, IsNil)
		c.Check(sec, HasLen, 1)

		buf := bytes.NewBuffer(nil)
		err = sec[0].Generate(buf)

		c.Check(err, IsNil)
		c.Check(buf.String(), Equals, test.expected, Commentf("Failed output check for test `%s`", test.name))
	}
}

func loadfile(filename string) string {
	buf := bytes.NewBuffer(nil)
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = io.Copy(buf, f)
	if err != nil {
		panic(err)
	}
	return string(buf.Bytes())
}
