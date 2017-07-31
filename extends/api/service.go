package api

import (
	"github.com/mbict/gogen"
)

type Service struct {
	gogen.Describer
	Name    string
	Methods []*Method
}

func (*Service) Context() string {
	return "api.service"
}

func (s *Service) AddMethod(m *Method) {
	s.Methods = append(s.Methods, m)
}

func NewService(name string) *Service {
	return &Service{
		Describer: &gogen.DescriptionExpr{},
		Name:      name,
	}
}
