package api

import (
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/lib"
)

type Service struct {
	gogen.Describer
	gogen.Packager
	Name         string
	Methods      []*Method
	HttpEndpoint *HTTP
}

type MethodIterator func(m *Method) error

func (*Service) Context() string {
	return "api.service"
}

func (s *Service) Finalize() {
	s.IterateMethods(func(m *Method) error {
		//add json tags to payload
		for _, field := range *gogen.AsObject(m.Payload.Type) {
			if _, ok := field.Attribute.Metadata["go:struct:tag:json"]; !ok {
				field.Attribute.Metadata["go:struct:tag:json"] = []string{lib.SnakeCase(field.Name), "omitempty"}
			}
		}

		//add json tags to response
		for _, field := range *gogen.AsObject(m.Results.Type) {
			if _, ok := field.Attribute.Metadata["go:struct:tag:json"]; !ok {
				field.Attribute.Metadata["go:struct:tag:json"] = []string{lib.SnakeCase(field.Name), "omitempty"}
			}
		}
		return nil
	})
}

func (s *Service) AddMethod(m *Method) {
	m.Service = s
	s.Methods = append(s.Methods, m)
}

func (s *Service) IterateMethods(it MethodIterator) error {
	//names := make([]string, len(s.Methods))
	//i := 0
	//for n := range r.Actions {
	//	names[i] = n
	//	i++
	//}
	//sort.Strings(names)
	//for _, n := range names {
	//if err := it(r.Actions[n]); err != nil {
	//		return err
	//	}
	//}
	for _, m := range s.Methods {
		it(m)
	}
	return nil
}

func NewService(name string) *Service {
	return &Service{
		Describer:    &gogen.DescriptionExpr{},
		Packager:     &gogen.PackageNameExpr{},
		Name:         name,
		HttpEndpoint: NewHTTP(),
	}
}
