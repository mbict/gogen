package codegen

import (
	"github.com/mbict/gogen/extends/api"
	"github.com/mbict/gogen/generator"
)

var gen *Generator

type Generator struct {
}

func init() {
	gen = &Generator{}
}

func Register() {
	generator.Register(gen)
}

func (g *Generator) Name() string {
	return "api"
}

func (g *Generator) Generate(path string) ([]generator.FileWriter, error) {
	cg := NewCodeGenerator(path)
	return cg.Writers(api.Root)
}
