package api

var Root *Api

type Api struct {
	Services []*Service
}

func (*Api) Context() string {
	return "api.Root"
}

func (a *Api) AddService(s *Service) {
	a.Services = append(a.Services, s)
}

func NewApi() *Api {
	return &Api{
	}
}


func init() {
	Root = NewApi()
}