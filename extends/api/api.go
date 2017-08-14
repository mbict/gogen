package api

import "github.com/mbict/gogen/dslengine"

var Root *Api

type Api struct {
	Services []*Service
}

func (a *Api) DSLName() string {
	return "API"
}

func (a *Api) IterateSets(iterator dslengine.SetIterator) {
	defs := make([]dslengine.Definition, len(a.Services))
	for i, s := range a.Services {
		defs[i] = s
	}
	iterator(defs)
}

func (a *Api) Reset() {
	panic("implement me")
}

func (*Api) Context() string {
	return "api.Root"
}

func (a *Api) AddService(s *Service) {
	a.Services = append(a.Services, s)
}

func NewApi() *Api {
	return &Api{}
}

func init() {
	Root = NewApi()

	dslengine.Register(Root)
}
