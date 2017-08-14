package codegen

import (
	"fmt"
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/generator"
)

type Generator struct {
}

type Codegen struct {
	generator.CodeGenerator
	basePackage string
}

var gen *Generator

func init() {
	gen = &Generator{}
	generator.Register(gen)
}

func Register() {
	generator.Register(gen)
}

func (g *Generator) Name() string {
	return "gogen"
}

func (g *Generator) Generate(path string) ([]generator.FileWriter, error) {
	//We set the absolute root package here
	gogen.Root.SetPackage("/" + path)

	//generate the codegenerator (if needed)
	//todo: check if this code generator is needed
	codegen := NewCodeGenerator(path)
	return codegen.Writers(gogen.Root)
}

func NewCodeGenerator(basePackage string) *Codegen {
	cg := generator.NewCodeGenerator()
	return &Codegen{
		CodeGenerator: cg,
		basePackage:   basePackage,
	}
}

func (g *Codegen) Writers(root interface{}) ([]generator.FileWriter, error) {
	root, ok := root.(*gogen.Gogen)
	if !ok {
		return nil, fmt.Errorf("Incompatible root")
	}

	res := []generator.FileWriter{}

	//do g3n here

	return res, nil
}

func (g *Codegen) GenerateType(ut *gogen.UserTypeExpr) ([]generator.Section, error) {
	//t := g.Template().Lookup("EVENT")
	//if t == nil {
	//	return nil, errors.New("template not found")
	//}
	//
	//imports := generator.NewImports(path.Join(g.basePackage, "domain/event"))
	//imports.Add("github.com/mbict/go-cqrs")
	//imports.AddFromAttribute(e.Attributes)
	//
	//s := generator.Section{
	//	Template: template.Must(t.Clone()),
	//	Data: map[string]interface{}{
	//		"Event":   e,
	//		"Imports": imports,
	//	},
	//}

	return []generator.Section{ /*s*/ }, nil
}
