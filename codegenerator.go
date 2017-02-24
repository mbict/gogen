package gogen

import (
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
	"github.com/mbict/gogen/lib"
)

type CodeGenerator interface {
	//Returns the basic template engine
	Template() *template.Template
}

type codeGenerator struct {
	template *template.Template
}

func NewCodeGenerator(templatePath ...string) CodeGenerator {

	t := template.New("base")

	t.Funcs(template.FuncMap{
		"snake":   lib.SnakeCase,
		"title":   strings.Title,
		"toLower": strings.ToLower,
		"toUpper": strings.ToUpper,
		"untitle": lib.UnTitle,
		//go specific functions
		"packageName": PackageName,
		"dict":        lib.Dictionary,
	})

	//base templates
	_, file, _, _ := runtime.Caller(0)
	pattern := filepath.Join(path.Dir(file), "templates", "*.tmpl")
	template.Must(t.ParseGlob(pattern))

	//user provided templates
	for _, path := range templatePath {
		template.Must(t.ParseGlob(path))
	}

	return &codeGenerator{
		template: t,
	}
}

func (c *codeGenerator) Template() *template.Template {
	return c.template
}
