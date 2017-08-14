package api

import (
	"github.com/mbict/gogen"
)

type Method struct {
	gogen.DescriptionExpr
	Service *Service
	Name    string
	Results *gogen.AttributeExpr
	Payload *gogen.AttributeExpr
	Errors  []*gogen.AttributeExpr

	//Endpoints used to access this service (HTTPEndpoint or GRPC)
	Endpoints []Endpoint
}

type EndpoiontIterator func(e Endpoint) error

func (*Method) Context() string {
	return "api.method"
}

func (m *Method) Finalize() {
}

func NewMethod(name string) *Method {
	return &Method{
		Name:    name,
		Results: gogen.NewAttribute(&gogen.Object{}),
		Payload: gogen.NewAttribute(&gogen.Object{}),
	}
}

func (m *Method) AddEndpoint(e Endpoint) {
	m.Endpoints = append(m.Endpoints, e)
}

func (m *Method) IterateEndpoints(it EndpoiontIterator) error {
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
	for _, e := range m.Endpoints {
		it(e)
	}
	return nil
}
