package codegen

import (
	"fmt"
	"github.com/mbict/gogen"
	"github.com/mbict/gogen/extends/api"
	"github.com/mbict/gogen/generator"
	"github.com/mbict/gogen/lib"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"text/template"
)

type Codegen struct {
	generator.CodeGenerator
	basePackage string
}

func NewCodeGenerator(basePackage string) *Codegen {
	_, file, _, _ := runtime.Caller(0)
	templatePath := filepath.Join(path.Dir(file), "templates", "*.tmpl")
	cg := generator.NewCodeGenerator(templatePath)
	return &Codegen{
		CodeGenerator: cg,
		basePackage:   basePackage,
	}
}

func (g *Codegen) Writers(root interface{}) ([]generator.FileWriter, error) {
	api, ok := root.(*api.Api)
	if !ok {
		return nil, fmt.Errorf("Incompatible root")
	}

	res := []generator.FileWriter{}
	for _, service := range api.Services {

		//-------------------------------------------
		// Services
		//-------------------------------------------
		log.Printf("Generating service interface for '%s'\n", service.Name)
		sections, err := g.GenerateServiceInterface(service)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("%s/%s.go", lib.SnakeCase(service.Name), "service")
		fwService := generator.NewFileWriter(sections, file)
		res = append(res, fwService)

		//-------------------------------------------
		// Endpoints
		//-------------------------------------------
		log.Printf("Generating endpoints for service '%s'\n", service.Name)
		sections, err = g.GenerateEndpointsInterface(service)
		if err != nil {
			return nil, err
		}

		file = fmt.Sprintf("%s/%s.go", lib.SnakeCase(service.Name), "endpoint")
		fwEndpoint := generator.NewFileWriter(sections, file)
		res = append(res, fwEndpoint)

		//-------------------------------------------
		// Transport HTTPEndpoint
		//-------------------------------------------
		log.Printf("Generating http transport layer for service '%s'\n", service.Name)
		sections, err = g.GenerateHTTPTransport(service)
		if err != nil {
			return nil, err
		}

		file = fmt.Sprintf("%s/%s.go", lib.SnakeCase(service.Name), "transport_http")
		fwHTTPTransport := generator.NewFileWriter(sections, file)
		res = append(res, fwHTTPTransport)
	}

	/*for _, a := range api.Aggregates {
		s, err := g.GenerateAggregate(a)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("api/aggregate/%s.go", lib.SnakeCase(a.Name))
		fwProjection := generator.NewFileWriter(s, file)
		res = append(res, fwProjection)

		for _, c := range a.Commands {
			s, err := g.GenerateCommand(c)
			if err != nil {
				return nil, err
			}

			file := fmt.Sprintf("api/command/%s.go", lib.SnakeCase(c.Name))
			fwCommand := generator.NewFileWriter(s, file)
			res = append(res, fwCommand)
		}
	}

	for _, e := range api.AllEvents() {
		s, err := g.GenerateServiceInterface(e)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("api/event/%s.go", lib.SnakeCase(e.Name))
		fwEvent := generator.NewFileWriter(s, file)
		res = append(res, fwEvent)
	}

	for _, p := range api.Projections {
		s, err := g.GenerateProjection(p)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("api/projection/%s.go", lib.SnakeCase(p.Name))
		fwProjection := generator.NewFileWriter(s, file)
		res = append(res, fwProjection)
	}

	s, err := g.GenerateAggregatesFactory(api)
	if err != nil {
		return nil, err
	}
	fwAggregateFactory := generator.NewFileWriter(s, "api/aggregate_factory.go")
	res = append(res, fwAggregateFactory)

	s, err = g.GenerateEventsFactory(api)
	if err != nil {
		return nil, err
	}
	fwEventFactory := generator.NewFileWriter(s, "api/event_factory.go")
	res = append(res, fwEventFactory)

	s, err = g.GenerateRepositoryInterfaces(api)
	if err != nil {
		return nil, err
	}
	fwRepository := generator.NewFileWriter(s, "repository/repository.go")
	res = append(res, fwRepository)

	for _, r := range api.AllReadRepositories() {
		s, err := g.GenerateDbRepository(r)
		if err != nil {
			return nil, err
		}

		file := fmt.Sprintf("repository/sql/%s_repository.go", lib.SnakeCase(r.Name))
		fwDBRepository := generator.NewFileWriter(s, file)
		res = append(res, fwDBRepository)
	}*/

	return res, nil
}

func (g *Codegen) GenerateServiceInterface(service *api.Service) ([]generator.Section, error) {

	imports := generator.NewImports(path.Join(service.Package()))
	t, _ := g.Template().Clone()
	t = t.Funcs(template.FuncMap{
		"package": func() string { return service.Package() },
	})

	imports.Add("context")
	//gather all the imports from the result types in the methods
	for _, method := range service.Methods {
		imports.AddFromAttribute(method.Results)
		imports.AddFromAttribute(method.Payload)
	}

	//create the user models
	sections := []generator.Section{}
	for _, ut := range gogen.Root.UserTypes {
		if ut.Package() != service.Package() {
			continue
		}

		imports.AddFromAttribute(ut.Attribute())

		sections = append(sections, generator.Section{
			Template: template.Must(t.Lookup("USERTYPE").Clone()),
			Data:     ut,
		})
	}

	//add package and imports header
	sections = append([]generator.Section{
		{
			Template: template.Must(t.Lookup("PACKAGE").Clone()),
			Data:     imports,
		},
	}, sections...)

	//add the service interface
	sections = append(sections,
		generator.Section{
			Template: template.Must(t.Lookup("SERVICE_INTERFACE").Clone()),
			Data: map[string]interface{}{
				"Service": service,
				"Imports": imports,
			},
		})

	return sections, nil
}

func (g *Codegen) GenerateEndpointsInterface(service *api.Service) ([]generator.Section, error) {
	imports := generator.NewImports(path.Join(service.Package()))
	t, _ := g.Template().Clone()
	t = t.Funcs(template.FuncMap{
		"package": func() string { return service.Package() },
	})

	//gather all the imports from the result types in the methods
	imports.Add("context")
	imports.Add("github.com/go-kit/kit/endpoint")
	//for _, m := range service.Methods {
	//	imports.AddFromAttribute(m.Results)
	//	imports.AddFromAttribute(m.Payload)
	//}

	//create the user models
	sections := []generator.Section{}
	//for _, ut := range gogen.Root.UserTypes {
	//	if ut.Package() != service.Package() {
	//		continue
	//	}
	//
	//	imports.AddFromAttribute(ut.Attribute())
	//
	//	sections = append(sections, generator.Section{
	//		Template: template.Must(t.Lookup("USERTYPE").Clone()),
	//		Data:     ut,
	//	})
	//}

	//add package and imports header
	sections = append([]generator.Section{
		{
			Template: template.Must(t.Lookup("PACKAGE").Clone()),
			Data:     imports,
		},
	}, sections...)

	//add the service interface
	sections = append(sections,
		generator.Section{
			Template: template.Must(t.Lookup("ENDPOINTS_INTERFACE").Clone()),
			Data: map[string]interface{}{
				"Service": service,
				"Imports": imports,
			},
		},
		generator.Section{
			Template: template.Must(t.Lookup("ENDPOINT_CONSTRUCTOR").Clone()),
			Data: map[string]interface{}{
				"Service": service,
				"Imports": imports,
			},
		})

	for _, method := range service.Methods {
		if len(method.Endpoints) == 0 {
			continue
		}

		imports.AddFromAttribute(method.Results)
		imports.AddFromAttribute(method.Payload)

		payloadObj := gogen.AsObject(method.Payload.Type)
		if !payloadObj.IsEmpty() {
			sections = append(sections,
				generator.Section{
					Template: template.Must(t.Lookup("ENDPOINT_REQUEST").Clone()),
					Data: map[string]interface{}{
						"Service": service,
						"Method":  method,
						"Payload": payloadObj,
						"Imports": imports,
					},
				})
		}

		sections = append(sections,
			generator.Section{
				Template: template.Must(t.Lookup("ENDPOINT_RESPONSE").Clone()),
				Data: map[string]interface{}{
					"Service": service,
					"Method":  method,
					"Results": payloadObj,
					"Imports": imports,
				},
			},
			generator.Section{
				Template: template.Must(t.Lookup("ENDPOINT_METHOD_CONSTRUCTOR").Clone()),
				Data: map[string]interface{}{
					"Service": service,
					"Method":  method,
					"Imports": imports,
				},
			})

	}

	return sections, nil
}

func (g *Codegen) GenerateHTTPTransport(service *api.Service) ([]generator.Section, error) {

	//extract the http endpoints
	var httpEndpoints []*api.HTTPEndpoint
	service.IterateMethods(func(m *api.Method) error {
		return m.IterateEndpoints(func(e api.Endpoint) error {
			if httpEndpoint, ok := e.(*api.HTTPEndpoint); ok {
				httpEndpoints = append(httpEndpoints, httpEndpoint)
			}
			return nil
		})
	})

	if len(httpEndpoints) == 0 {
		return nil, nil
	}

	imports := generator.NewImports(path.Join(service.Package()))
	t, _ := g.Template().Clone()
	t = t.Funcs(template.FuncMap{
		"package": func() string { return service.Package() },
	})

	//gather all the imports from the result types in the methods
	imports.Add("context")
	imports.Add("encoding/json")
	imports.Add("github.com/go-kit/kit/log")
	imports.AddAlias("httptransport", "github.com/go-kit/kit/transport/http")
	imports.AddAlias("mux", "github.com/mbict/httprouter")
	imports.Add("net/http")
	//for _, m := range service.Methods {
	//	imports.AddFromAttribute(m.Results)
	//	imports.AddFromAttribute(m.Payload)
	//}

	//create the user models

	//for _, ut := range gogen.Root.UserTypes {
	//	if ut.Package() != service.Package() {
	//		continue
	//	}
	//
	//	imports.AddFromAttribute(ut.Attribute())
	//
	//	sections = append(sections, generator.Section{
	//		Template: template.Must(t.Lookup("USERTYPE").Clone()),
	//		Data:     ut,
	//	})
	//}

	sections := []generator.Section{}

	//add package and imports header
	sections = append([]generator.Section{
		{
			Template: template.Must(t.Lookup("PACKAGE").Clone()),
			Data:     imports,
		},
	}, sections...)

	sections = append(sections,
		generator.Section{
			// @todo: need to make template.Must(t.Lookup().Clone) nil save
			Template: template.Must(t.Lookup("TRANSPORT_HTTP_CONSTRUCTOR").Clone()),
			Data: map[string]interface{}{
				"Service":   service,
				"Endpoints": httpEndpoints,
			},
		})

	////add the service interface
	//sections = append(sections,
	//	generator.Section{
	//		Template: template.Must(t.Lookup("ENDPOINTS_INTERFACE").Clone()),
	//		Data: map[string]interface{}{
	//			"Service": service,
	//			"Imports": imports,
	//		},
	//	},
	//	generator.Section{
	//		Template: template.Must(t.Lookup("ENDPOINT_CONSTRUCTOR").Clone()),
	//		Data: map[string]interface{}{
	//			"Service": service,
	//			"Imports": imports,
	//		},
	//	})

	//for _, method := range service.Methods {
	//	if len(method.Endpoints) == 0 {
	//		continue
	//	}
	//
	//	imports.AddFromAttribute(method.Results)
	//	imports.AddFromAttribute(method.Payload)
	//
	//	payloadObj := gogen.AsObject(method.Payload.Type)
	//	if !payloadObj.IsEmpty() {
	//		sections = append(sections,
	//			generator.Section{
	//				Template: template.Must(t.Lookup("ENDPOINT_REQUEST").Clone()),
	//				Data: map[string]interface{}{
	//					"Service": service,
	//					"Method":  method,
	//					"Payload": payloadObj,
	//					"Imports": imports,
	//				},
	//			})
	//	}
	//
	//	sections = append(sections,
	//	)
	//
	//	sections = append(sections,
	//		generator.Section{
	//			Template: template.Must(t.Lookup("ENDPOINT_RESPONSE").Clone()),
	//			Data: map[string]interface{}{
	//				"Service": service,
	//				"Method":  method,
	//				"Results": payloadObj,
	//				"Imports": imports,
	//			},
	//		},
	//		generator.Section{
	//			Template: template.Must(t.Lookup("ENDPOINT_METHOD_CONSTRUCTOR").Clone()),
	//			Data: map[string]interface{}{
	//				"Service": service,
	//				"Method":  method,
	//				"Imports": imports,
	//			},
	//		})
	//
	//}

	return sections, nil
}

/*
func (g *Codegen) GenerateCommand(c *cqrs.CommandExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("COMMAND")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "domain/command"))
	imports.Add(gogen.UUID.Package)
	imports.AddFromAttribute(c.Params)

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Command": c,
			"Imports": imports,
		},
	}

	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateAggregate(a *cqrs.AggregateExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("AGGREGATE")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "domain/aggregate"))
	imports.Add("errors")
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "domain/event"))
	imports.Add(path.Join(g.basePackage, "domain/command"))

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Aggregate": a,
			"Imports":   imports,
		},
	}

	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateProjection(p *cqrs.ProjectionExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("PROJECTION")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "domain/projection"))
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(path.Join(g.basePackage, "domain/event"))

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Projection": p,
			"Imports":    imports,
		},
	}

	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateAggregatesFactory(d *cqrs.DomainExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("AGGREGATE_FACTORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "domain"))
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "domain/aggregate"))

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}

	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateEventsFactory(d *cqrs.DomainExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("EVENT_FACTORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "domain"))
	imports.Add("github.com/mbict/go-cqrs")
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "domain/event"))

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}

	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateRepositoryInterfaces(d *cqrs.DomainExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("REPOSITORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "repository"))
	imports.Add(gogen.UUID.Package)
	imports.Add(path.Join(g.basePackage, "models"))
	for _, r := range d.AllReadRepositories() {
		if r.Filter != nil {
			imports.AddFromAttribute(r.Filter)
		}
	}

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Domain":  d,
			"Imports": imports,
		},
	}
	return []generator.Section{s}, nil
}

func (g *Codegen) GenerateDbRepository(r *cqrs.RepositoryExpr) ([]generator.Section, error) {
	t := g.Template().Lookup("DB_REPOSITORY")
	if t == nil {
		return nil, errors.New("template not found")
	}

	imports := generator.NewImports(path.Join(g.basePackage, "repository/sql"))
	imports.Add("database/sql")
	imports.Add(gogen.UUID.Package)
	imports.Add("github.com/masterminds/squirrel")
	imports.Add(path.Join(g.basePackage, "repository"))
	imports.Add(path.Join(g.basePackage, "models"))
	if r.Filter != nil {
		imports.AddFromAttribute(r.Filter)
	}

	s := generator.Section{
		Template: template.Must(t.Clone()),
		Data: map[string]interface{}{
			"Repository": r,
			"Imports":    imports,
		},
	}

	return []generator.Section{s}, nil
}*/
