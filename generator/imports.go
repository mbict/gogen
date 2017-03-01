package generator

import (
	"github.com/mbict/gogen/lib"
	"path"
	"regexp"
	"github.com/mbict/gogen"
)

type Imports interface {
	Add(string)
	AddFromAttribute(*gogen.AttributeExpr)
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

func (i *imports) Add(packagePath string) {
	i.packages[packagePath] = packagePath
}

// AddFromAttribute collects all the external packages needed by the attribute definition
func (i *imports) AddFromAttribute(attr *gogen.AttributeExpr) {
	if attr == nil {
		return
	}

	var recursiveImport func(interface{})
	recursiveImport = func(in interface{}) {
		switch t := in.(type) {
		case gogen.Composite:
			recursiveImport(t.Attribute().Type)
		case *gogen.UserTypeExpr:
			if t.Package != "" {
				i.packages[t.Package] = t.Package
			}
		case *gogen.Array:
			recursiveImport(t.ElemType.Type)
		case *gogen.Object:
			for _, field := range *t {
				recursiveImport(field.Attribute.Type)
			}
		}
	}
	recursiveImport(attr.Type)
}

func (i *imports) ToSlice() []string {
	return lib.StringMapToSlice(i.packages)
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
