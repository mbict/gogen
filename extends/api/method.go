package api

import (
	"github.com/mbict/gogen"
)

type Method struct {
	gogen.Describer
	Name   string
	Result *gogen.AttributeExpr
}

func (*Method) Context() string {
	return "api.method"
}

func NewMethod(name string) *Method {
	return &Method{
		Describer: &gogen.DescriptionExpr{},
		Name:      name,
	}
}
