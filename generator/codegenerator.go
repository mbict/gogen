package generator

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/lib"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"
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
		"snake":     lib.SnakeCase,
		"camel":     func(str string) string { return lib.VarName(str, true) },
		"title":     strings.Title,
		"lowercase": strings.ToLower,
		"uppercase": strings.ToUpper,
		"untitle":   lib.UnTitle,
		"join":      strings.Join,
		"quote":     strconv.Quote,
		"unquote":   strconv.Unquote,
		"padleft":   lib.PadLeft,
		"padright":  lib.PadRight,

		//attribute datatype specific functions
		"isObject":    gogen.IsObject,
		"isArray":     gogen.IsArray,
		"isPrimitive": gogen.IsPrimitive,
		"isMap":       gogen.IsMap,

		//go specific functions
		"varname":     lib.VarName,
		"packageName": PackageName,
		"dict":        lib.Dictionary,
		"baseName":    path.Base,
		"package":     func() string { return "" }, //empty stub
	})

	//base templates
	_, file, _, _ := runtime.Caller(0)
	pattern := filepath.Join(path.Dir(file), "/../codegen/templates", "*.tmpl")
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
