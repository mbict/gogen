package generator

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/lib"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type Imports interface {
	Add(pkg string)
	AddAlias(alias string, pkg string)
	AddFromAttribute(*gogen.AttributeExpr)
	Package() string
	ToSlice() []string
}

type imports struct {
	packages    map[string]string
	basePackage string
}

func NewImports(basePackage string) Imports {
	return &imports{
		packages:    make(map[string]string, 0),
		basePackage: basePackage,
	}
}

//type TypeToImportMap map[DataType]string
//
//var DefaultTypeMapping = TypeToImportMap{
//	UUID: "github.com/satori/go.uuid",
//}

func (i *imports) Add(pkg string) {
	i.packages[pkg] = pkg
}

// AddAlias will add the new package or overwrite an existing package
func (i *imports) AddAlias(alias string, pkg string) {
	i.packages[pkg] = alias
}

// AddFromAttribute collects all the external packages needed by the attribute definition
func (i *imports) AddFromAttribute(attr *gogen.AttributeExpr) {
	if attr == nil {
		return
	}

	var recursiveImport func(interface{})
	recursiveImport = func(in interface{}) {
		switch t := in.(type) {
		case gogen.Packager:
			if t.Package() != "" && t.Package() != i.basePackage {
				i.packages[t.Package()] = t.Package()
			}
		case *gogen.Array:
			recursiveImport(t.ElemType.Type)
		case gogen.Composite:
			recursiveImport(t.Attribute().Type)
		case *gogen.Object:
			for _, field := range *t {
				recursiveImport(field.Attribute.Type)
			}
		}
	}
	recursiveImport(attr.Type)
}

func (i *imports) ToSlice() []string {
	mapKeys := lib.MapIndex(i.packages)
	res := make([]string, len(i.packages))
	for index, key := range mapKeys {
		if !strings.EqualFold(i.packages[key], key) {
			//aliased
			res[index] = i.packages[key] + " " + strconv.Quote(key)
		} else {
			res[index] = strconv.Quote(key)
		}
	}
	return res
}

func (i *imports) Package() string {
	return i.basePackage
}

var basePackageRegex = regexp.MustCompile(`^(?:go[.\-_])?(.*)$`)

func PackageName(importPath string) string {
	name := path.Base(importPath)
	res := basePackageRegex.FindAllStringSubmatch(name, 1)
	if len(res) >= 1 && len(res[0]) >= 1 {
		return res[0][1]
	}
	return name
}
